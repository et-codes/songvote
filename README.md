# SongVote (in development)

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
- `PUT /songs/{id}` updates song information
- `DELETE /songs/{id}` remove a song from the list
- `GET /songs/vote/{id}` get votes for a song
- `POST /songs/vote` add vote for a song
- `POST /songs/veto` add veto for a song

### Users (DONE)

- `POST /users` adds a user
- `GET /users/{id}` returns user information
- `PUT /users/{id}` updates user information
- `DELETE /users/{id}` deletes a user

## TODO

- Add logging to server errors and successes
- Allow for undoing a vote or a veto? (TBD)
