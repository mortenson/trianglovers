package main

func getStrings() map[string][]string {
	return map[string][]string{
		"COMFORT_A": []string{"You're sitting in front of the TV - what chair are you in?"},
		"COMFORT_A_STRONG": []string{
			"Recliner, easy. Feet up, drink in my coaster, remote in hand.",
			"Anywhere is fine, as long as I don't have to wear pants.",
		},
		"COMFORT_A_NORMAL": []string{
			"Maybe a nice sofa I can sink into.",
			"I like to curl up in my love seat like a cat.",
		},
		"COMFORT_A_DEFAULT": []string{
			"I don't watch TV, and prefer to stand.",
			"I'm on the ground, doing push ups every time the studio\naudience laughs.",
		},
		"WEALTH_A": []string{"What's your dream birthday gift?"},
		"WEALTH_A_STRONG": []string{
			"An iced out wrist watch is always a good choice.",
			"Can't have too many designer bags.",
		},
		"WEALTH_A_NORMAL": []string{
			"A dinner out at the best restaurant in town.",
			"Maybe a wine of good vintage.",
		},
		"WEALTH_A_DEFAULT": []string{
			"Always into a good book or a gift card.",
			"I love me some knick knacks.",
		},
		"ADVENTURE_A": []string{"What trip is on your bucket list?"},
		"ADVENTURE_A_STRONG": []string{
			"I would go to the bottom of the ocean.",
			"Climbing a mountain a third time sounds fun.",
		},
		"ADVENTURE_A_NORMAL": []string{
			"I've always wanted to eat street food from another country.",
			"There's an island nearby I'd love to sail to.",
		},
		"ADVENTURE_A_DEFAULT": []string{
			"I'm more of a homebody, most of my trips are over my cats!.",
			"I took a trip to our local import goods store once.",
		},
		"EXCITEMENT_A": []string{"What really gets your heart racing?"},
		"EXCITEMENT_A_STRONG": []string{
			"Doing 30 on my unicycle, bombing hills with no helment on.",
			"Sky diving without a parachute - I've tried doing it without a partner\nbut they wouldn't let me.",
		},
		"EXCITEMENT_A_NORMAL": []string{
			"Paintball. Getting pelted in the face is a rush.",
			"I'm in an arm wrestling league that can get pretty nuts.",
		},
		"EXCITEMENT_A_DEFAULT": []string{
			"Ordering delivery when I know I have food in the fridge.",
			"Fixing a dropped stitch while I'm knitting.",
		},
		"ROMANCE_A": []string{"Describe your ideal date."},
		"ROMANCE_A_STRONG": []string{
			"I cook an incredible meal, rose petals lead to the bedroom, I'm\nsprawled out on the bed covered in baby oil.",
			"I read you my poetry in front of a roaring fire as we melt into\na fur rug.",
		},
		"ROMANCE_A_NORMAL": []string{
			"A mutual foot massage always lights my fire.",
			"Just a simple cuddle is all I need.",
		},
		"ROMANCE_A_DEFAULT": []string{
			"A piping hot TV-dinner and a good re-run.",
			"Going to see my team play while they get me beer.",
		},
		"FAMILY_A": []string{"How do you feel about kids?"},
		"FAMILY_A_STRONG": []string{
			"I already have names and schools picked out!",
			"They're our future. My future, Our future.",
		},
		"FAMILY_A_NORMAL": []string{
			"I just want to pinch their little points!",
			"I definitely want them, just give me a few years first!",
		},
		"FAMILY_A_DEFAULT": []string{
			"They're just tiny, smelly, stupid versions of adults.",
			"I could take or leave them. Mostly leave.",
		},
	}
}
