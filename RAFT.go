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
func supportTheKing(support []chan int, petition chan int, haveAKing *bool){
	var canidate int;
	for {
		select {
		case canidate = <-petition:
			if (!(*haveAKing)){
				support[canidate] <- 1;
				*haveAKing = true;
			}
		}
	}
}

//propose yourself as king and listen for the response
func doItYourself(i int, support, petition []chan int, haveAKing *bool){
	var votesGarnered int
	cycleTicker := time.NewTicker(CYCLE_FREQUENCY * time.Millisecond);
	for {
		select {
		case <-cycleTicker.C:
			*haveAKing = false;
			initiative := MIN_WAIT + (CYCLE_FREQUENCY - MIN_WAIT) * rand.Float64()
			time.Sleep(time.Duration(initiative) * time.Second)
			if (!(*haveAKing)){
				for j:=0; j<NUM_NODES; j++{
					if (i!=j){
						petition[j] <- i
					}
				}
				*haveAKing = true
			}
			votesGarnered = 0
		case <-support[i]:
			votesGarnered += 1
			if (votesGarnered > NUM_NODES/2 + 1){
				fmt.Println("I'm the king", i)
			}
		}
	}
}

//Run the raft election algorithm
func RAFT(i int, done chan bool, wg *sync.WaitGroup, petition, support []chan int){
	fmt.Println("start ", i);
	haveAKing := false;
	go supportTheKing(support, petition[i], &haveAKing);
	go doItYourself(i, support, petition, &haveAKing);

	//wait on shutdown
	for {
		select {
		case <-done:
			wg.Done()
			return
		}
	}
}

func main() {
	petition := make([]chan int, NUM_NODES); 
	for i:=0; i<NUM_NODES; i++{
		petition[i] = make(chan int);
	}
	support := make([]chan int, NUM_NODES); 
	for i:=0; i<NUM_NODES; i++{
		support[i] = make(chan int);
	}
	
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(NUM_NODES)
	for i:=0; i<NUM_NODES; i++{
		go RAFT(i, done, &wg, petition, support);
	}

	time.Sleep(RUN_TIME * time.Second)
	for i:=0; i<NUM_NODES; i++{
		done <- true
	}
	wg.Wait()

}