// Implementação do algoritmo de travessia de Tarry
package main

import (
	"fmt"
	"time"
)

type Message struct {
	Sender string
	IdMax int
}

type Neighbour struct {
	Id   string
	From chan Message
	To   chan Message
}

func redirect(in chan Message, neigh Neighbour) {
	message := <-neigh.From
	in <- message
}

func process(id string, myid int, message Message, neighs ...Neighbour) {
	var pai Neighbour

	// Redeirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Message, 1)
	nmap := make(map[string]Neighbour)
	for _, neigh := range neighs {
		nmap[neigh.Id] = neigh
		go redirect(in, neigh)
	}

	if message.Sender == "init" {
		// Processo iniciador
		fmt.Printf("* %s é inicializador.\n", id)
		message.Sender = id
		message.IdMax= myid
		neighs[0].To <- message
		size := len(neighs)

		for i := 1; i < size; i++ {
			message := <-in
			fmt.Printf("From %s to %s\n", message.Sender, id)
			message.Sender = id
			message.IdMax = myid
			neighs[i].To <- message
		}

		message := <-in
		fmt.Printf("From %s to %s\n", message.Sender, id)
		fmt.Println("Fim!")
	} else {
		// Processo não iniciador
		tk := <-in
		//fmt.Printf("From %s to %s\n", tk.Sender, id)

		for _, neigh := range neighs {
			if pai.Id == "" {
				pai = nmap[tk.Sender]
				if (myid<tk.IdMax){
					myid=tk.IdMax
				}
				fmt.Printf("Explorando o processo %s \n", id)
			}

		// Entrega o message para o vizinho se ele não for o pai
		if pai.Id != neigh.Id {
			tk.Sender = id
			tk.IdMax = myid
			neigh.To <- tk
			tk = <-in
			fmt.Printf("From %s to %s\n", tk.Sender, id)
		}
		}

		// Token volta para o pai depois de ter passado enviado para todos os vizinhos
		tk.Sender = id
		tk.IdMax = myid
		pai.To <- tk

		fmt.Printf("Process: %s Eleito: %v", id, myid)
		
	}

}

func main() {

	pW := make(chan Message, 1)
	pS := make(chan Message, 1)
	pR := make(chan Message, 1)
	wP := make(chan Message, 1)
	wS := make(chan Message, 1)
	sP := make(chan Message, 1)
	sW := make(chan Message, 1)
	rQ := make(chan Message, 1)
	rP := make(chan Message, 1)
	qR := make(chan Message, 1)

	
	go process("W",1, Message{}, Neighbour{"S", sW, wS}, Neighbour{"P", pW, wP})
	go process("S",2, Message{}, Neighbour{"W", wS, sW},Neighbour{"P", pS, sP})
	go process("R",3, Message{}, Neighbour{"Q", qR, rQ}, Neighbour{"P", pR, rP})
	go process("Q",4, Message{}, Neighbour{"R", rQ, qR})
	go process("P",5, Message{"init",5}, Neighbour{"W", wP, pW}, Neighbour{"S", sP, pS}, Neighbour{"R", rP, pR})
	time.Sleep(5 * time.Second);
	fmt.Printf("* Fim da execução do algoritmo.")
}
