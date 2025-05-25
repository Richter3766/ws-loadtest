package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"ws-loadtest/metrics"
	"ws-loadtest/wsclient"
)

var (
    wsURL string
    roomID  int
    numClients int
    numRepeats int

)

func init() {
    flag.StringVar(&wsURL, "wsURL", "ws://localhost:8080/api/ws/websocket", "웹소켓 서버 주소")
    flag.IntVar(&roomID, "roomID", 1, "방 ID")
    flag.IntVar(&numClients, "numClients", 300, "동시 접속자 수")
    flag.IntVar(&numRepeats, "numRepeats", 1, "반복 수")
}

func main() {
    flag.Parse()
    metricsDir := fmt.Sprintf("result/test_%dclients_%drepeats", numClients, numRepeats)
    if err := os.MkdirAll(metricsDir, os.ModePerm); err != nil {
        log.Fatalf("폴더 생성 실패: %v", err)
    }

    for round := 1; round <= numRepeats; round++ {
        log.Printf("==== %d번째 부하 테스트 시작 ====", round)
        var wg sync.WaitGroup

        stopAll := make(chan struct{})
        users, err := wsclient.LoadUsers("./users.json")
        if err != nil {
            log.Fatalf("유저 로드 실패: %v", err)
        }

        // 지표 리셋
        metrics.ResetAll()

        var readyWg sync.WaitGroup
        readyWg.Add(numClients)
        start := make(chan struct{})

        // 클라이언트 고루틴 기동
        for i := 1; i <= numClients; i++ {
            wg.Add(1)
            jwt := users[i-1].JWT
            go wsclient.RunClient(i, wsURL, roomID, jwt, stopAll, &wg,  &readyWg, start)
        }
        readyWg.Wait()
        close(start)
        // 회차별 부하테스트: 일정 시간 대기 후 종료 신호
        time.Sleep(300 * time.Second)
 
        close(stopAll)
        log.Println("종료 신호 감지: 모든 클라이언트 종료 대기")
        wg.Wait()

        // 회차별 metrics 저장
        fname := fmt.Sprintf("%s/metrics_%d.json", metricsDir, round)
        metrics.SaveMetricsJSON(fname)
        log.Printf("==== %d번째 테스트 종료, 결과 저장: %s ====", round, fname)

        // 충분한 대기시간(자원 해제) 후 다음 라운드(필요시)
        time.Sleep(10 * time.Second)
    }
    log.Println("모든 부하 테스트 완료")
}