# KLADS

## Description  
This project is a web platform designed to integrate essential academic and communication tools for students and instructors. It combines features inspired by Canvas (Learning Management System) and Campuswire (class communication tool) to provide a seamless user experience.  

## I AM ONLY ADDING THINGS THAT I HAVE WORKED ON IN THIS README FILE. ## 

Key features include:  
- **Login and Authentication**: Secure sign-up, sign-in, and sign-out functionalities.  
- **Databases**: Robust database schema for storing users, topics, comments, assignments, and messages.  
- **Direct Messaging (DMs)**: Private, secure communication between users.  
- **Chatbot**: AI-powered chatbot using OpenAI's API to assist users with quick answers to their queries.  

## Features  

### 1. **Login and Authentication**  
- Secure credential management using bcrypt for password hashing.  
- JSON Web Tokens (JWT) for session management and authentication.  
- Role-based access control for students and instructors.  

### 2. **Databases**  
The platform utilizes a relational database to manage and store essential data:  
- **Users**: Stores user credentials, roles, and profile information.  
- **Topics**: Maintains forum discussions and threads.  
- **Comments**: Tracks comments and replies on forum topics.  
- **Assignments**: Handles assignment uploads and metadata (e.g., deadlines, file links).  
- **Messages**: Logs all direct messages for secure and private communication.  

### 3. **Direct Messaging (DMs)**  
- Real-time private messaging between users.  
- Ensures message security through encryption.  
- Optimized for minimal latency using WebSocket technology.  

### 4. **Chatbot**  
- Integrated with OpenAI’s API to provide intelligent, contextual responses.  
- Handles free-text input from users and fetches relevant information from the database.  
- Supports multi-turn conversations to enhance user experience.  
- Common use cases include answering questions about assignments, deadlines, and platform usage.  

## Setup Instructions  

### Prerequisites  
- [Go](https://golang.org/) installed on your system.  
- Database (PostgreSQL recommended).  
- API key for OpenAI’s GPT-4 (if using the chatbot functionality).  

### Steps  
1. Clone the repository:  

   git clone git@github.com:MattLADS/website.git
   cd website
   git checkout KpathaK

4. Run the project:  

   go run .
  

5. Access the application:  
   Open your browser and go to `http://localhost:8080`.   

## Future Enhancements  
- Add video conferencing integration.  
- Implement advanced analytics for assignments and messaging.  
- Extend chatbot capabilities to include sentiment analysis and real-time feedback.  
 

## Acknowledgments  
- OpenAI for providing the GPT-4 API.  
- Go community for the robust frameworks and tools.  
- Team members and contributors who made this project possible.  
  
