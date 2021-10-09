package user

var (
	SampleUsers = map[int]User{
		User1.Uid: User1,
		User2.Uid: User2,
	}
)

var (
	User1 = User{
		Uid:    1,
		Name:   "zouying",
		Avatar: "https://nathanleclaire.com/images/iowriter/aviator.png",
	}

	User2 = User{
		Uid:    2,
		Name:   "gopher",
		Avatar: "http://img01.yohoboys.com/contentimg/2018/11/22/13/0187be5a52edcdc999f749b9e24c7815fb.jpg",
	}
)
