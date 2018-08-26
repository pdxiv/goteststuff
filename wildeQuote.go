package main

// A fun little project to auto-generate Oscar Wilde quotes using Markov Chains

// - One refinement idea, is that if there is no more than 1 option for
// choosing the next word, with the current lookbehind setting, go to a lower
// lookbehind setting, to make sure that things are more interesting. This
// should make expressions both more grammatically correct, while keeping things
// interesting.

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const setSize = 2 // Number of token combinations to get stats for

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	wildeQuote := [...]string{
		"A cynic is a man who knows the price of everything, and the value of nothing.",
		"A gentleman is one who never hurts anyone's feelings unintentionally.",
		"A good friend will always stab you in the front.",
		"A little sincerity is a dangerous thing, and a great deal of it is absolutely fatal.",
		"A man can be happy with any woman, as long as he does not love her.",
		"A man can't be too careful in the choice of his enemies.",
		"A man who does not think for himself does not think at all.",
		"A man's face is his autobiography. A woman's face is her work of fiction.",
		"A thing is not necessarily true because a man dies for it.",
		"A work of art is the unique result of a unique temperament.",
		"All women become like their mothers. That is their tragedy. No man does, and that is his.",
		"Always forgive your enemies; nothing annoys them so much.",
		"Ambition is the germ from which all growth of nobleness proceeds.",
		"Ambition is the last refuge of the failure.",
		"America had often been discovered before Columbus, but it had always been hushed up.",
		"America is the only country that went from barbarism to decadence without civilization in between.",
		"An idea that is not dangerous is unworthy of being called an idea at all.",
		"Anybody can be good in the country. There are no temptations there.",
		"Anybody can sympathise with the sufferings of a friend, but it requires a very fine nature to sympathise with a friend's success.",
		"Anyone who lives within their means suffers from a lack of imagination.",
		"Arguments are extremely vulgar, for everyone in good society holds exactly the same opinion.",
		"Arguments are to be avoided: they are always vulgar and often convincing.",
		"Art is individualism, and individualism is a disturbing and disintegrating force.",
		"Art is the most intense mode of individualism that the world has known.",
		"Art never harms itself by keeping aloof from the social problems of the day: rather, by so doing, it more completely realises for us that which we desire.",
		"Art should never try to be popular. The public should try to make itself artistic.",
		"As long as war is regarded as wicked, it will always have its fascination. When it is looked upon as vulgar, it will cease to be popular.",
		"Be yourself; everyone else is already taken.",
		"Behind every exquisite thing that existed, there was something tragic.",
		"Between men and women there is no friendship possible. There is passion, enmity, worship, love, but no friendship.",
		"Bigamy is having one wife too many. Monogamy is the same.",
		"By giving us the opinions of the uneducated, journalism keeps us in touch with the ignorance of the community.",
		"Children begin by loving their parents; after a time they judge them; rarely, if ever, do they forgive them.",
		"Consistency is the last refuge of the unimaginative.",
		"Conversation about the weather is the last refuge of the unimaginative.",
		"Death and vulgarity are the only two facts in the nineteenth century that one cannot explain away.",
		"Deceiving others. That is what the world calls a romance.",
		"Democracy means simply the bludgeoning of the people by the people for the people.",
		"Education is an admirable thing, but it is well to remember from time to time that nothing that is worth knowing can be taught.",
		"Every portrait that is painted with feeling is a portrait of the artist, not of the sitter.",
		"Every saint has a past, and every sinner has a future.",
		"Everybody who is incapable of learning has taken to teaching.",
		"Everything in moderation, including moderation.",
		"Everything in the world is about sex except sex. Sex is about power.",
		"Everything popular is wrong.",
		"Experience is merely the name men gave to their mistakes.",
		"Experience is one thing you can't get for nothing.",
		"Fashion is a form of ugliness so intolerable that we have to alter it every six months.",
		"Fathers should be neither seen nor heard. That is the only proper basis for family life.",
		"Hatred is blind, as well as love.",
		"He has no enemies, but is intensely disliked by his friends.",
		"How can a woman be expected to be happy with a man who insists on treating her as if she were a perfectly normal human being.",
		"How marriage ruins a man! It is as demoralizing as cigarettes, and far more expensive.",
		"I always pass on good advice. It is the only thing to do with it. It is never of any use to oneself.",
		"I am not young enough to know everything.",
		"I am so clever that sometimes I don't understand a single word of what I am saying.",
		"I am the only person in the world I should like to know thoroughly.",
		"I am too fond of reading books to care to write them.",
		"I can resist everything except temptation.",
		"I can stand brute force, but brute reason is quite unbearable. There is something unfair about its use. It is hitting below the intellect.",
		"I choose my friends for their good looks, my acquaintances for their good characters, and my enemies for their intellects. A man cannot be too careful in the choice of his enemies.",
		"I don't want to be at the mercy of my emotions. I want to use them, to enjoy them, and to dominate them.",
		"I don't want to go to heaven. None of my friends are there.",
		"I find it harder and harder every day to live up to my blue china.",
		"I have never given adoration to any body except myself.",
		"I have nothing to declare except my genius.",
		"I have the simplest tastes. I am always satisfied with the best.",
		"I like persons better than principles, and I like persons with no principles better than anything else in the world.",
		"I never travel without my diary. One should always have something sensational to read in the train.",
		"I put all my genius into my life; I put only my talent into my works.",
		"I see when men love women. They give them but a little of their lives. But women when they love give everything.",
		"I sometimes think that God in creating man somewhat overestimated his ability.",
		"I suppose society is wonderfully delightful. To be in it is merely a bore. But to be out of it is simply a tragedy.",
		"I think God, in creating man, somewhat overestimated his ability.",
		"I want my food dead. Not sick, not dying, dead.",
		"If one cannot enjoy reading a book over and over again, there is no use in reading it at all.",
		"If one plays good music, people don't listen and if one plays bad music people don't talk.",
		"If there was less sympathy in the world, there would be less trouble in the world.",
		"If you are not too long, I will wait here for you all my life.",
		"If you pretend to be good, the world takes you very seriously. If you pretend to be bad, it doesn't. Such is the astounding stupidity of optimism.",
		"If you want to tell people the truth, make them laugh, otherwise they'll kill you.",
		"In all matters of opinion, our adversaries are insane.",
		"In America the President reigns for four years, and Journalism governs forever and ever.",
		"In America the young are always ready to give to those who are older than themselves the full benefits of their inexperience.",
		"In married life three is company and two none.",
		"In modern life nothing produces such an effect as a good platitude. It makes the whole world kin.",
		"It is a very sad thing that nowadays there is so little useless information.",
		"It is absurd to divide people into good and bad. People are either charming or tedious.",
		"It is better to be beautiful than to be good. But, it is better to be good than to be ugly.",
		"It is better to have a permanent income than to be fascinating.",
		"It is only an auctioneer who can equally and impartially admire all schools of art.",
		"It is only by not paying one's bills that one can hope to live in the memory of the commercial classes.",
		"It is only the modern that ever becomes old-fashioned.",
		"It is the spectator, and not life, that art really mirrors.",
		"It is through art, and through art only, that we can realise our perfection.",
		"It is what you read when you don't have to that determines what you will be when you can't help it.",
		"It takes great deal of courage to see the world in all its tainted glory, and still to love it.",
		"Keep love in your heart. A life without it is like a sunless garden when the flowers are dead.",
		"Laughter is not at all a bad beginning for a friendship, and it is far the best ending for one.",
		"Let us have no machine-made ornament at all; it is all bad and worthless and ugly.",
		"Life imitates art far more than art imitates Life.",
		"Life is far too important a thing ever to talk seriously about.",
		"Life is never fair, and perhaps it is a good thing for most of us that it is not.",
		"Live! Live the wonderful life that is in you! Let nothing be lost upon you. Be always searching for new sensations. Be afraid of nothing.",
		"Man can believe the impossible, but man can never believe the improbable.",
		"Man is a rational animal who always loses his temper when he is called upon to act in accordance with the dictates of reason.",
		"Man is least himself when he talks in his own person. Give him a mask, and he will tell you the truth.",
		"Memory, is the diary that we all carry about with us.",
		"Men always want to be a woman's first love; women like to be a man's last romance.",
		"Men marry because they are tired; women, because they are curious; both are disappointed.",
		"Moderation is a fatal thing. Nothing succeeds like excess.",
		"Morality is simply the attitude we adopt towards people whom we personally dislike.",
		"Most people are other people. Their thoughts are someone else's opinions, their lives a mimicry, their passions a quotation.",
		"Most people die of a sort of creeping common sense, and discover when it is too late that the only things one never regrets are one's mistakes.",
		"Never love anyone who treats you like you're ordinary.",
		"No good deed goes unpunished.",
		"No great artist ever sees things as they really are. If he did, he would cease to be an artist.",
		"No man is rich enough to buy back his past.",
		"No object is so beautiful that, under certain conditions, it will not look ugly.",
		"No woman should ever be quite accurate about her age. It looks so calculating.",
		"Nothing can cure the soul but the senses, just as nothing can cure the senses but the soul.",
		"Nothing is so aggravating than calmness.",
		"Now that the House of Commons is trying to become useful, it does a great deal of harm.",
		"Nowadays people know the price of everything and the value of nothing.",
		"One can survive everything, nowadays, except death, and live down everything except a good reputation.",
		"One of the many lessons that one learns in prison is, that things are what they are and will be what they will be.",
		"One should always be in love. That is the reason one should never marry.",
		"One should always play fairly when one has the winning cards.",
		"One should never trust a woman who tells one her real age. A woman who would tell one that would tell one anything.",
		"One's past is what one is. It is the only way by which people should be judged.",
		"One's real life is so often the life that one does not lead.",
		"Only dull people are brilliant at breakfast.",
		"Only the shallow know themselves.",
		"Ordinary riches can be stolen; real riches cannot. In your soul are infinitely precious things that cannot be taken from you.",
		"Our ambition should be to rule ourselves, the true kingdom for each one of us; and true progress is to know more, and be more, and to do more.",
		"Patriotism is the virtue of the vicious.",
		"Perhaps, after all, America never has been discovered. I myself would say that it had merely been detected.",
		"Pessimist: One who, when he has the choice of two evils, chooses both.",
		"Questions are never indiscreet, answers sometimes are.",
		"Quotation is a serviceable substitute for wit.",
		"Ridicule is the tribute paid to the genius by the mediocrities.",
		"Romance should never begin with sentiment. It should begin with science and end with a settlement.",
		"Selfishness is not living as one wishes to live, it is asking others to live as one wishes to live.",
		"Seriousness is the only refuge of the shallow.",
		"Society exists only as a mental concept; in the real world there are only individuals.",
		"Society often forgives the criminal; it never forgives the dreamer.",
		"Some cause happiness wherever they go; others whenever they go.",
		"Success is a science; if you have the conditions, you get the result.",
		"The advantage of the emotions is that they lead us astray.",
		"The basis of optimism is sheer terror.",
		"The books that the world calls immoral are books that show the world its own shame.",
		"The critic has to educate the public; the artist has to educate the critic.",
		"The difference between literature and journalism is that journalism is unreadable and literature is not read.",
		"The function of the artist is to invent, not to chronicle.",
		"The good ended happily, and the bad unhappily. That is what fiction means.",
		"The heart was made to be broken.",
		"The imagination imitates. It is the critical spirit that creates.",
		"The man who can dominate a London dinner-table can dominate the world.",
		"The moment you think you understand a great work of art, it's dead for you.",
		"The mystery of love is greater than the mystery of death.",
		"The nicest feeling in the world is to do a good deed anonymously; and have somebody find out.",
		"The old believe everything, the middle-aged suspect everything, the young know everything.",
		"The one charm about marriage is that it makes a life of deception absolutely necessary for both parties.",
		"The only difference between the saint and the sinner is that every saint has a past, and every sinner has a future.",
		"The only good thing to do with good advice is pass it on; it is never of any use to oneself.",
		"The only thing to do with good advice is to pass it on. It is never of any use to oneself.",
		"The only way to get rid of temptation is to yield to it. I can resist everything but temptation.",
		"The public is wonderfully tolerant. It forgives everything except genius.",
		"The salesman knows nothing of what he is selling save that he is charging a great deal too much for it.",
		"The spirit of an age may be best expressed in the abstract ideal arts, for the spirit itself is abstract and ideal.",
		"The true mystery of the world is the visible, not the invisible.",
		"The truth is rarely pure and never simple.",
		"The very essence of romance is uncertainty.",
		"The well bred contradict other people. The wise contradict themselves.",
		"The world has grown suspicious of anything that looks like a happily married life.",
		"The world is a stage, but the play is badly cast.",
		"The world is divided into two classes, those who believe the incredible, and those who do the improbable.",
		"There are many things that we would throw away if we were not afraid that others might pick them up.",
		"There are only two kinds of people who are really fascinating; people who know absolutely everything, and people who know absolutely nothing.",
		"There are only two tragedies in life: one is not getting what one wants, and the other is getting it.",
		"There is a luxury in self-reproach. When we blame ourselves we feel no one else has a right to blame us.",
		"There is always something infinitely mean about other people's tragedies.",
		"There is always something ridiculous about the emotions of people whom one has ceased to love.",
		"There is no necessity to separate the monarch from the mob; all authority is equally bad.",
		"There is no sin except stupidity.",
		"There is no such thing as a moral or an immoral book. Books are well written, or badly written.",
		"There is nothing in the world like the devotion of a married woman. It is a thing no married man knows anything about.",
		"There is nothing so difficult to marry as a large nose.",
		"There is only one class in the community that thinks more about money than the rich, and that is the poor. The poor can think of nothing else.",
		"There is only one thing in life worse than being talked about, and that is not being talked about.",
		"There is something terribly morbid in the modern sympathy with pain. One should sympathise with the colour, the beauty, the joy of life. The less said about life's sores the better.",
		"There's nothing in the world like the devotion of a married woman. It's a thing no married man knows anything about.",
		"This suspense is terrible. I hope it will last.",
		"Those who find ugly meanings in beautiful things are corrupt without being charming. This is a fault.",
		"Those whom the gods love grow young.",
		"To define is to limit.",
		"To expect the unexpected shows a thoroughly modern intellect.",
		"To live is the rarest thing in the world. Most people exist, that is all.",
		"To lose one parent may be regarded as a misfortune; to lose both looks like carelessness.",
		"To love oneself is the beginning of a lifelong romance.",
		"True friends stab you in the front.",
		"We are all in the gutter, but some of us are looking at the stars.",
		"We are each our own devil, and we make this world our hell.",
		"We live in an age when unnecessary things are our only necessities.",
		"What is a cynic? A man who knows the price of everything and the value of nothing.",
		"What we have to do, what at any rate it is our duty to do, is to revive the old art of Lying.",
		"When a man has once loved a woman he will do anything for her except continue to love her.",
		"When I was young I thought that money was the most important thing in life; now that I am old I know that it is.",
		"When the gods wish to punish us they answer our prayers.",
		"Whenever a man does a thoroughly stupid thing, it is always from the noblest motives.",
		"Whenever people agree with me I always feel I must be wrong.",
		"Who, being loved, is poor?",
		"With freedom, books, flowers, and the moon, who could not be happy?",
		"Woman begins by resisting a man's advances and ends by blocking his retreat.",
		"Women are made to be loved, not understood.",
		"Women are never disarmed by compliments. Men always are. That is the difference between the sexes.",
		"Women love us for our defects. If we have enough of them, they will forgive us everything, even our gigantic intellects.",
		"Work is the curse of the drinking classes.",
		"You can never be overdressed or overeducated.",
	}

	var token []string
	// Pattern for splitting out non-word characters
	re := regexp.MustCompile("^([a-zA-Z'-]*)([^a-zA-Z'-]*)$")
	// Take inventory of which tokens exist. Sort out duplicates with a map.
	tokenMap := make(map[string]int)

	for _, quote := range wildeQuote {
		lowerCaseQuote := strings.ToLower(quote)
		spaceSeparatedQuote := strings.Split(lowerCaseQuote, " ")
		for _, v := range spaceSeparatedQuote {
			submatch := re.FindStringSubmatch(v)
			tokenMap[submatch[1]] = 0
			if submatch[2] != "" {
				tokenMap[submatch[2]] = 0
			}
		}
	}

	// Put tokens in an array to get a fixed order
	for k := range tokenMap {
		token = append(token, k)
	}

	var stats = map[[setSize]int]map[int]int{}
	var statsOccurrance = map[[setSize]int]int{}

	for _, quote := range wildeQuote {
		lowerCaseQuote := strings.ToLower(quote)
		spaceSeparatedQuote := strings.Split(lowerCaseQuote, " ")
		// Convert the quote string into a list of tokens
		var tokenTextList []string
		for _, v := range spaceSeparatedQuote {
			submatch := re.FindStringSubmatch(v)
			tokenTextList = append(tokenTextList, submatch[1])
			if submatch[2] != "" {
				tokenTextList = append(tokenTextList, submatch[2])
			}
		}

		// Convert the text tokens into numerical codes
		var tokenList []int
		// Add a "start of line" token(s) at beginning
		for i := 0; i < setSize; i++ {
			tokenList = append(tokenList, -1)
		}
		for _, textToken := range tokenTextList {
			for i, v := range token {
				if v == textToken {
					tokenList = append(tokenList, i)
				}
			}
		}
		// Add an "end of line token"
		tokenList = append(tokenList, -2)

		// Gather statistics
		// Details for doing this: https://stackoverflow.com/questions/44305617/nested-maps-in-golang
		var derp [setSize]int
		for i := 0; i < len(tokenList)-setSize; i++ {
			// Map of maps stuff below
			for j := 0; j < setSize; j++ {
				derp[j] = tokenList[i+j]
			}
			// Check if map is initialized. Otherwise, initialize it
			if stats[derp] == nil {
				stats[derp] = map[int]int{}
			}
			stats[derp][tokenList[i+setSize]]++
		}
	}

	for k, v := range stats {
		for _, occurrences := range v {
			statsOccurrance[k] += occurrences // To be used for randomizing
		}
	}

	var generatedList []int
	// Initialize generated slice with "start of sentence" markers
	for i := 0; i < setSize; i++ {
		generatedList = append(generatedList, -1)
	}
	var derp [setSize]int
	for generatedList[len(generatedList)-1] != -2 {
		// Fill array that is used for map key
		for i := 0; i < setSize; i++ {
			derp[i] = generatedList[len(generatedList)-setSize+i]
		}
		// Select a random value out of the possibilities
		randomSelection := rand.Intn(statsOccurrance[derp])
		counter := 0
		for k, v := range stats[derp] {
			counter += v
			if randomSelection < counter {
				generatedList = append(generatedList, k)
				if k == -2 {
					break
				}
				break
			}
		}
	}
	generatedList = generatedList[setSize:]
	generatedList = generatedList[0 : len(generatedList)-1]
	for _, v := range generatedList {
		fmt.Print(token[v], " ")
	}
	fmt.Println()
}
