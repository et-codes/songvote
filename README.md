# SongVote

## Under development...

**SongVote** is an app to allow users to:
- Adding songs, voting, and vetoing are done in rounds
- Add songs to a list (users have limited number of songs to add per round)
- Vote for the songs in the list
- Veto a song on the list (users have limited number of vetos per round)
- Use vote results to generate a list of "approved" songs for the round

## API
### Songs
- `GET /songs` returns a list of all songs
- `POST /songs` adds a song to the list
- `GET /songs/{id}` info for particular song
- `PATCH /songs/{id}` updates song information
- `DELETE /songs/{id}` remove a song from the list
- `POST /songs/vote/{id}` add vote for a song
- `POST /songs/veto/{id}` add veto for a song

### Users
- `POST /users` adds a user
- `GET /users/{id}` returns user information
- `PATCH /users/{id}` updates user information
- `DELETE /users/{id}` deletes a user

### TODO
- Make song's ID its own type alias for int64
- Add Song PATCH
- Add Song POST vote
- Add Song POST veto
- Allow for undoing a vote or a veto
- Send JSON error messages from server.go
- Add log output
