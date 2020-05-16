package main

func getStrings() map[string][]string {
	return map[string][]string{
		"COMFORT_A": []string{"You're sitting in front of the TV - what chair are you in?"},
		"COMFORT_A_STRONG": []string{
			"Recliner, easy. Feet up, drink in my coaster, remote in hand.",
			"Anywhere is fine, as long as I don't have to wear pants.",
			"A bean bag chair so big I literally can't get out of it.",
		},
		"COMFORT_A_NORMAL": []string{
			"Maybe a nice sofa I can sink into.",
			"I like to curl up in my love seat like a cat.",
			"My recliner wearing a sweater twice my size.",
		},
		"COMFORT_A_DEFAULT": []string{
			"I don't watch TV, and prefer to stand.",
			"I'm on the ground, doing push ups every time the studio\naudience laughs.",
			"I usually watch shows while pacing my house, pretending\nI'm part of the cast.",
		},
		"WEALTH_A": []string{"What's your dream birthday gift?"},
		"WEALTH_A_STRONG": []string{
			"An iced out wrist watch is always a good choice.",
			"Can't have too many designer bags.",
			"Something from the heart...like a third yacht.",
		},
		"WEALTH_A_NORMAL": []string{
			"A dinner out at the best restaurant in town.",
			"Maybe a wine of good vintage.",
			"A box of fine imported cigars.",
		},
		"WEALTH_A_DEFAULT": []string{
			"Always into a good book or a gift card.",
			"I love me some knick knacks!",
			"Anything homemade is perfect.",
		},
		"ADVENTURE_A": []string{"What trip is on your bucket list?"},
		"ADVENTURE_A_STRONG": []string{
			"Diving to the ocean floor to woo a merperson. That is, unless\nthis works out.",
			"Climbing another mountain sounds fun.",
			"Hitchhiking wherever the road takes me.",
		},
		"ADVENTURE_A_NORMAL": []string{
			"I've always wanted to eat street food from another country.",
			"There's an island nearby I'd love to sail to.",
			"Going deep into the wilderness to connect with nature.",
		},
		"ADVENTURE_A_DEFAULT": []string{
			"I'm more of a homebody, most of my trips are over my cats!.",
			"I took a trip to our local import goods store once.",
			"My character in the MMO I'm playing does all the adventuring I need.",
		},
		"EXCITEMENT_A": []string{"What really gets your heart racing?"},
		"EXCITEMENT_A_STRONG": []string{
			"Doing 30 on my unicycle, bombing hills with no helmet on.",
			"Sky diving without a parachute - I've tried doing it without a partner\nbut they wouldn't let me.",
			"I'm a magician's assistant, you know the kind they throw real\nknives at?",
		},
		"EXCITEMENT_A_NORMAL": []string{
			"Paintball. Getting pelted in the face is a rush.",
			"I'm in an arm wrestling league that can get pretty nuts.",
			"My friends and I still parkour, it's cool until you get hurt,then it's\nSUPER cool.",
		},
		"EXCITEMENT_A_DEFAULT": []string{
			"Ordering delivery when I know I have food in the fridge.",
			"Fixing a dropped stitch while I'm knitting.",
			"Buying shoes without trying them on first.",
		},
		"ROMANCE_A": []string{"Describe your ideal date."},
		"ROMANCE_A_STRONG": []string{
			"I cook an incredible meal, rose petals lead to the bedroom, I'm\nsprawled out on the bed covered in baby oil.",
			"I read you my poetry in front of a roaring fire as we melt into\na fur rug.",
			"We share a fondue fork, we share a bottle of wine, then we\nshare our deepest desires.",
		},
		"ROMANCE_A_NORMAL": []string{
			"A mutual foot massage always lights my fire.",
			"Just a simple cuddle is all I need.",
			"When they wipe sauce off my face with no judgement in their heart.",
		},
		"ROMANCE_A_DEFAULT": []string{
			"A piping hot TV-dinner and a good re-run.",
			"Watching my team play while my partner gets me beer.",
			"We have an early meal, then I fall asleep on the couch while they\nclean up.",
		},
		"FAMILY_A": []string{"How do you feel about kids?"},
		"FAMILY_A_STRONG": []string{
			"I already have names and schools picked out!",
			"They're our future. My future. Our future.",
			"I already have four and still want more. Hah!",
		},
		"FAMILY_A_NORMAL": []string{
			"I just want to pinch their little points!",
			"I definitely want them, just give me a few years first!",
			"They're like pets that learn how to talk. What's not to love?",
		},
		"FAMILY_A_DEFAULT": []string{
			"They're just tiny, smelly, stupid versions of adults.",
			"I could take or leave them. Mostly leave.",
			"I wouldn't want to have me as a parent.",
		},
	}
}
