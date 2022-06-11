package handler

import (
	"SecondHandMarketBackend/model"
	"SecondHandMarketBackend/service"
	"fmt"
	"net/http"
)

func messageNewUploadHandler(w http.ResponseWriter, r *http.Request) {
	//new conversation for two user
	//known attrs: user1id, user2id, init message(string)
	/*need to do:
	1:get two users via orm first.
	2:make a new conversation including init message for this two users 
	and orm automatically update the table
	3:return to the front end
	*/
	user1_id := r.Context().Value("user1_id")
	user2_id := r.Context().Value("user2_id")
	init_message := r.Context().Value("message")
	message := model.Message{
		Message:    init_message.(string),
		SenderId:   user1_id.(int),
		ReceiverId: user2_id.(int),
	}
	s := []model.Message{message}
	conversation := model.Conversation{
		//conversation need two model.User object to be user1 and user2
		//todo
		User1Id:  user1_id.(int),
		User2Id:  user2_id.(int),
		Messages: s,
	}
	err := service.CreateConversation(&conversation)
	if err != nil {
		http.Error(w, "Failed to save data to backend", http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "New Conversation established")
}
