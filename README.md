# SongVote

## Under development...

**SongVote** is an app to allow users to:
- Adding songs, voting, and vetoing are done in rounds
- Add songs to a list (users have limited number of songs to add per round)
- Vote for the songs in the list
- Veto a song on the list (users have limited number of vetos per round)
- Use vote results to generate a list of "approved" songs for the round
- When a round ends, the song list resets. The song list from previous rounds is stored. Vetoes are resupplied to the users.

## API
### Songs (DONE)
- `GET /songs` returns a list of all songs
- `POST /songs` adds a song to the list
- `GET /songs/{id}` info for particular song
- `PATCH /songs/{id}` updates song information
- `DELETE /songs/{id}` remove a song from the list
- `POST /songs/vote/{id}` add vote for a song
- `POST /songs/veto/{id}` add veto for a song

### Users (WIP)
- `POST /users` adds a user
- `GET /users/{id}` returns user information
- `PATCH /users/{id}` updates user information
- `DELETE /users/{id}` deletes a user

### TODO
- Allow for undoing a vote or a veto? (TBD)
- Implement users
  - Get, add, update, and delete
  - Track who adds a song
  - Track who voted for and vetoed a song
  - Password encryption
