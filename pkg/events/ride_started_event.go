package events

type RideStartedMessage struct {
	DriverId string `json:"driver_id"`
	RideId   string `json:"rider_id"`
}

type DriverLocationEvent struct {
	DriverId string  `json:"driver_id"`
	RiderId  string  `json:"rider_id"`
	Lat      float32 `json:"lat"`
	Long     float32 `json:"long"`
}
