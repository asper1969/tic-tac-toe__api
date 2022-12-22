# Welcome to Tic-tac-toe skills API 

App based on gobuffalo web framework

## Database Setup

It looks like you chose to set up your application using a database! Fantastic!

The first thing you need to do is open up the "database.yml" file and edit it to use the correct usernames, passwords, hosts, etc... that are appropriate for your environment.

You will also need to make sure that **you** start/install the database of your choice. Buffalo **won't** install and start it for you.

### Create Your Databases

Ok, so you've edited the "database.yml" file and started your database, now Buffalo can create the databases in that file for you:

```console
buffalo pop create -a
```

## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

```console
buffalo dev
```

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.

**Congratulations!** You now have your Buffalo application up and running.

## Starting socket worker
```console
buffalo task socket_server:run 
```

## Room notification statuses
1 - Notification is ready for sending

2 - Notification was sent


## Update session request types (room notification types)
1 - Accept move event initiated by Mr.Moderator

2 - Question answer event initiated by team

3 - Answer window created

## TODO:

Add token verification for sessions

Migrate to Postgre. Change datetime fields type in models. Implement postgre pubsub for notifications

## Tournaments:
Moderator creates session. In respons he gets session code

Each Team join session by session code

Moderator initiates tournament start

At that point teams see their opponents and can start game

Moderator gets realtime statistics from matches

Moderator can stop tournament at any time

Modearator and teams can relaunch the app and join back to tournament by unique tournament token

