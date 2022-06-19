package domain

type Statistics struct {
	ElectionId string `json:"election_id"`
	Age        int    `json:"age" bson:"age"`
	Region     string `json:"region" bson:"region"`
	Circuit    string `json:"circuit" bson:"circuit"`
	Sex        string `json:"sex" bson:"sex"`
	Capacity   int    `json:"capacity" bson:"capacity"`
	Votes      int    `json:"votes" bson:"votes"`
}
