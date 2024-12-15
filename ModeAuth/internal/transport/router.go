package transport

import (
	"ModeAuth/internal/service"
	"ModeAuth/pkg/logging"
	"log"
	"net/http"
)

func RunRouter(aService service.IAuth, sService service.IStates, addr string) {
	c := NewController(aService, sService)

	http.HandleFunc("/user/auth/check", c.CheckUser())
	http.HandleFunc("/user/auth/check-block", c.CheckUserIsBlocked())
	http.HandleFunc("/user/state/go-work", c.GoWorkState())
	http.HandleFunc("/user/state/sending", c.SendingState())
	http.HandleFunc("/user/state/checking", c.CheckingState())
	http.HandleFunc("/user/state/get-report", c.GetReport())

	log.Printf(logging.INFO+"API is running on port %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(logging.FATAL+"Failed to start server:", err)
	}
}
