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
- `POST /songs` add a song to the list
- `GET /songs/{song}` info for particular song
- `POST /songs/{song}` add vote or veto for a song
- `PATCH /songs/{song}` updates song information
- `DELETE /songs/{song}` remove a song from the list

### Users
- `POST /users` adds a user
- `GET /users/{user}` returns user information
- `DELETE /users/{user}` deletes a user
