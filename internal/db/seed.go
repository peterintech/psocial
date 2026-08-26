package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/peterintech/psocial/internal/store"
)

var usernames = []string{
	"John", "Jane", "Michael", "Emily", "William", "Olivia", "James", "Ava", "Benjamin", "Sophia", "Daniel", "Isabella", "Matthew", "Mia", "Joseph", "Charlotte", "David", "Amelia", "Andrew", "Evelyn", "Christopher", "Harper", "Joshua", "Luna", "Nathan", "Scarlett", "Ryan", "Victoria", "Samuel", "Aurora", "Alexander", "Penelope", "Ethan", "Nora", "Jacob", "Hannah", "Mason", "Aria", "Logan", "Layla", "Elijah", "Lucas", "Abigail", "Oliver", "Elizabeth", "Liam", "Sofia", "Aiden", "Ella", "Noah", "Madison", "Caden", "Lily", "Jackson", "Avery", "Grayson", "Zoe", "Henry", "Grace", "Alexander", "Chloe", "Isaac", "Victoria", "Jack",
	"Lily",
	"Owen",
	"Hannah",
	"Wyatt",
	"Addison",
	"Sebastian",
	"Natalie",
	"Jayden",
	"Samantha",
	"Carter",
	"Leah",
	"Evan",
	"Sarah",
	"Luke",
	"Megan",
	"Kevin",
	"Rachel",
	"Brian",
	"Jessica",
	"Edward",
	"Ashley",
	"Ronald",
	"Emily",
	"Timothy",
	"Diane",
	"Jason",
	"Julie",
	"Jeffrey",
	"Joyce",
	"Ryan",
	"Evelyn",
	"Jacob",
	"Joan",
	"Gary",
	"Victoria",
	"Nicholas",
	"Kelly",
	"Eric",
	"Christina",
	"Jonathan",
	"Lauren",
}

var titles = []string{
	"Exploring the Wonders of Nature",
	"Unveiling the Secrets of the Universe",
	"The Art of Mindfulness: Finding Inner Peace",
	"Journey to the Depths of the Ocean",
	"Mastering the Craft of Photography",
	"Discovering Hidden Gems in Your City",
	"The Science Behind Happiness and Well-being",
	"Adventures in the World of Culinary Delights",
	"Unraveling the Mysteries of Ancient Civilizations",
	"The Power of Positive Thinking: Transform Your Life",
	"Embracing Change: Navigating Life's Transitions",
	"Exploring the Intersection of Technology and Society",
	"The Beauty of Minimalism: Simplifying Your Life",
	"Unlocking Your Creative Potential: Tips and Techniques",
	"The Journey of Self-Discovery: Finding Your True Purpose",
	"Exploring the Wonders of Space: A Cosmic Adventure",
	"The Art of Storytelling: Captivating Your Audience",
	"Discovering the Healing Power of Nature",
	"The Science of Sleep: Unlocking the Secrets to Restful Nights",
	"Embracing Diversity: Celebrating Differences in Our World",
}

var contents = []string{
	"Great post! I really enjoyed reading it.",
	"Thanks for sharing this information. It was very helpful.",
	"I completely agree with your points. Well said!",
	"This is a thought-provoking article. It made me think.",
	"I appreciate the effort you put into writing this. Excellent work!",
	"Interesting perspective! I hadn't considered that before.",
	"Your writing style is engaging and easy to follow.",
	"This topic is very relevant and important. Thank you for addressing it.",
	"I learned something new from your post. Keep up the good work!",
	"This is a well-researched piece. I can tell you put a lot of effort into it.",
	"Very informative! I'll definitely be referring to this in the future.",
	"Your insights are valuable and thought-provoking. Thank you for sharing.",
	"This is a well-written article. I appreciate the clarity and organization.",
	"Your post has inspired me to take action. Thank you for motivating me.",
	"I found your examples and anecdotes very relatable. Great job!",
	"This is a comprehensive guide on the topic. I learned a lot from it.",
	"Your post has sparked an interesting discussion in the comments section.",
	"I appreciate the practical tips you provided. They will be useful in my own work.",
	"This is a well-balanced analysis of the subject. I appreciate your objectivity.",
	"Your post has encouraged me to explore this topic further. Thank you for the inspiration.",
}

func Seed(store store.Storage, db *sql.DB) error {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	tx.Commit()

	posts := generatePosts(users, 200)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			return fmt.Errorf("failed to create post: %w", err)
		}
	}

	comments := generateComments(posts, 500)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}
	}

	log.Println("Database seeding completed successfully.")
	return nil
}

func generatePosts(users []*store.User, n int) []*store.Post {
	posts := make([]*store.Post, n)
	for i := range n {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags:    []string{"tag1", "tag2"},
			Version: 1,
		}
	}
	return posts
}

func generateComments(posts []*store.Post, n int) []*store.Comment {
	comments := make([]*store.Comment, n)
	for i := range n {
		post := posts[rand.Intn(len(posts))]
		comments[i] = &store.Comment{
			Content: contents[rand.Intn(len(contents))],
			PostID:  post.ID,
			UserID:  post.UserID,
		}
	}

	return comments
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)
	for i := range n {
		users[i] = &store.User{
			Username: usernames[rand.Intn(len(usernames))] + fmt.Sprintf("%d", i),
			Email:    usernames[rand.Intn(len(usernames))] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}
	return users
}
