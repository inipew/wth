package main

import (
	"context"
	"fmt"
	"flag"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	statsService "stats/v2ray-core/app/stats/command"
	"stats/v2ray-core/common/units"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	apiServerAddrPtr string
	apiTimeout       int
	apiJSON          bool
)

func main() {
	cmd := flag.NewFlagSet("stats", flag.ExitOnError)
	executeStats(cmd, os.Args[1:])
}

func SetSharedFlags(cmd *flag.FlagSet) {
	setSharedFlags(cmd)
}

func setSharedFlags(cmd *flag.FlagSet) {
	cmd.StringVar(&apiServerAddrPtr, "s", "127.0.0.1:8383", "API server address")
	cmd.StringVar(&apiServerAddrPtr, "server", "127.0.0.1:8383", "API server address")
	cmd.IntVar(&apiTimeout, "t", 3, "API timeout in seconds")
	cmd.IntVar(&apiTimeout, "timeout", 3, "API timeout in seconds")
	cmd.BoolVar(&apiJSON, "json", false, "Output in JSON format")
}

func executeStats(cmd *flag.FlagSet, args []string) {
	SetSharedFlags(cmd)
	var (
		runtime bool
		regexp  bool
		reset   bool
	)
	cmd.BoolVar(&runtime, "runtime", false, "Get runtime statistics")
	cmd.BoolVar(&regexp, "regexp", false, "Use regular expressions for filtering")
	cmd.BoolVar(&reset, "reset", false, "Reset statistics after fetching")
	cmd.Parse(args)
	unnamed := cmd.Args()

	if runtime {
		getRuntimeStats(apiJSON)
		return
	}

	// Execute the stats fetching logic
	client, conn := setupGrpcClient()
	defer conn.Close()

	stats := fetchStats(client, unnamed, regexp, reset)
	showStats(stats)
}

func setupGrpcClient() (statsService.StatsServiceClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient(apiServerAddrPtr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := statsService.NewStatsServiceClient(conn)
	return client, conn
}

func fetchStats(client statsService.StatsServiceClient, patterns []string, regexp, reset bool) []*statsService.Stat {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(apiTimeout)*time.Second)
	defer cancel()

	req := &statsService.QueryStatsRequest{
		Patterns: patterns,
		Reset_:   reset,
	}

	queryResp, err := client.QueryStats(ctx, req)
	if err != nil {
		log.Fatalf("could not query stats: %v", err)
	}

	sort.Slice(queryResp.Stat, func(i, j int) bool {
		return queryResp.Stat[i].Name < queryResp.Stat[j].Name
	})

	return queryResp.Stat
}

func showStats(stats []*statsService.Stat) {
    if len(stats) == 0 {
        return
    }
    formats := []string{"%-12s", "%s"}
    sum := int64(0)
    sb := new(strings.Builder)
    idx := 0
    writeRow(sb, 0, 0,
        []string{"Value", "Name"},
        formats,
    )
    for _, stat := range stats {
        idx++
        sum += stat.Value
        writeRow(
            sb, 0, idx,
            []string{units.ByteSize(stat.Value).String(), stat.Name},
            formats,
        )
    }
    sb.WriteString(
        fmt.Sprintf("\nTotal: %s\n", units.ByteSize(sum)),
    )
    os.Stdout.WriteString(sb.String())
}

func writeRow(sb *strings.Builder, indent, index int, values, formats []string) {
    if index == 0 {
        sb.WriteString(strings.Repeat(" ", indent+4))
    } else {
        sb.WriteString(fmt.Sprintf("%s%-4d", strings.Repeat(" ", indent), index))
    }
    for i, v := range values {
        format := "%-14s"
        if i < len(formats) {
            format = formats[i]
        }
        sb.WriteString(fmt.Sprintf(format, v))
    }
    sb.WriteByte('\n')
}

func getRuntimeStats(jsonOutput bool) {
    client, conn := setupGrpcClient()
	defer conn.Close()

    r := &statsService.SysStatsRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(apiTimeout)*time.Second)
	defer cancel()

    resp, err := client.GetSysStats(ctx, r)
    if err != nil {
        log.Fatalf("failed to get sys stats: %v", err)
    }
    if jsonOutput {
        showJSONResponse(resp)
        return
    }
    showRuntimeStats(resp)
}

func showRuntimeStats(s *statsService.SysStatsResponse) {
    formats := []string{"%-22s", "%-10s"}
    rows := [][]string{
        {"Up time", (time.Duration(s.Uptime) * time.Second).String()},
        {"Memory obtained", units.ByteSize(s.Sys).String()},
        {"Number of goroutines", fmt.Sprintf("%d", s.NumGoroutine)},
        {"Heap allocated", units.ByteSize(s.Alloc).String()},
        {"Live objects", fmt.Sprintf("%d", s.LiveObjects)},
        {"Heap allocated total", units.ByteSize(s.TotalAlloc).String()},
        {"Heap allocate count", fmt.Sprintf("%d", s.Mallocs)},
        {"Heap free count", fmt.Sprintf("%d", s.Frees)},
        {"Number of GC", fmt.Sprintf("%d", s.NumGC)},
        {"Time of GC pause", (time.Duration(s.PauseTotalNs) * time.Nanosecond).String()},
    }
    sb := new(strings.Builder)
    writeRow(sb, 0, 0,
        []string{"Item", "Value"},
        formats,
    )
    for i, r := range rows {
        writeRow(sb, 0, i+1, r, formats)
    }
    os.Stdout.WriteString(sb.String())
}

func protoToJSONString(m proto.Message, prefix, indent string) (string, error) { // nolint: unparam
	return strings.TrimSpace(protojson.MarshalOptions{Indent: indent}.Format(m)), nil
}

func showJSONResponse(m proto.Message) {
	output, err := protoToJSONString(m, "", "")
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", m)
		log.Fatalf("error encode json: %s", err)
	}
	fmt.Println(output)
}