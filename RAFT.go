package main
import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

type SkuttleBut struct {
	hotStory uint
	story Skuttle
}

type Skuttle struct{
	node_num uint
	heartBeat uint
	responseChan chan bool
	lastSeen time.Time
}

const RUN_TIME = 100
const CYCLE_FREQUENCY = 300
const MIN_WAIT = 150
const NUM_NODES = 8;

//Listen and respond to the first canidate in a cycle
func supportTheKing(i int, votes []chan int, haveAKing *bool, cycleTicker Ticker){
	var canidate int;
	for {
		select {
		case canidate<-votes[i]:
			if (!(*haveAKing)){
				support[canidate] <- 1;
				*haveAKing = 1;
			}
		}
	}
}

//propose yourself as king and listen for the response
func doItYourself(i int, votes []chan int, haveAKing *bool, cycleTicker Ticker){
	var votesGarnered int
	for {
		select {
		case <-cycleTicker.C:
			*haveAKing = false;
			initiative := MIN_WAIT + (CYCLE_FREQUENCY - MIN_WAIT) * rand.Float64()
			time.Sleep(initiative * time.Second)
			if (!(*haveAKing)){
				for j:=0; j<NUM_NODES; j++{
					if (i!=j){
						petition[j] <- i
					}
				}
				*haveAKing = true
			}
			votesGarnered = 0
		}
		case <- support[i]:
			votesGarnered += 1
			if (votesGarnered > NUM_NODES/2 + 1){
				//we are the king?
			}
		}
	}
}

//Run the raft election algorithm
func RAFT(i int, done chan bool, wg *sync.WaitGroup, petition []chan int, support []chan int){
	fmt.Println("start ", i);
	haveAKing := false;
	ticker := time.NewTicker(CYCLE_FREQUENCY * time.Millisecond);
	go supportTheKing(i, votes, &haveAKing, ticker);
	go doItYourself(i, votes, &haveAKing, ticker);

	//reset the king every cycle
	for {
		select {
		case <-done:
			wg.Done()
			return
		}
	}
}

func main() {
	messages := make([]chan Skuttle, NUM_NODES); 
	for i:=0; i<NUM_NODES; i++{
		messages[i] = make(chan Skuttle);
	}
	
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(NUM_NODES)
	for i:=0; i<NUM_NODES; i++{
		go RAFT(i, done, &wg, messages);
	}

	time.Sleep(RUN_TIME * time.Second)
	for i:=0; i<NUM_NODES; i++{
		done <- true
	}
	wg.Wait()

}