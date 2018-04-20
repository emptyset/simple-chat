# Domain

*Users*
- id
- username
- hash 

*Messages*
- id
- timestamp
- sender_id
- recipient_id
- content
- metadata

## Examples
```
user: { 
	id: 1,
	username: "emptyset",
	hash: "098j2f09j..."	// hashed using bcrypt library
}

message: {
	id: 10,
    timestamp: 1524258753,	// UTC, Unix
	sender_id: 1, 			// FK to user
	recipient_id: 2,		// FK to user
	content: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	// marshal'd metadata map[string][string]
    metadata: "{ _type: \"video\", length: \"72h3m0.5s\", source: \"YouTube\" }" 
}
```

# API
```
GET  /user/:id
PUT  /user/:id
  { username: "emptyset", password: "correct horse battery staple" }
POST /user
  { username: "emptyset", password: "incorrect horse battery staple" }

GET  /message?sid=1&rid=2&c=10&o=2
POST /message
  { sender_id: 1, recipient_id: 2, content: "...", metadata: "{ _type: \"video\", length: \"...\", ... }" }
```


