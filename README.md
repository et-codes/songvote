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
- `POST /song` adds a song to the list
- `GET /song/{id}` info for particular song
- `POST /song/{id}` add vote or veto for a song
- `PATCH /song/{id}` updates song information
- `DELETE /song/{id}` remove a song from the list

### Users
- `POST /users` adds a user
- `GET /users/{id}` returns user information
- `DELETE /users/{id}` deletes a user

### TODO
- Make song's ID its own type alias for int64
- Add the missing Song methods
