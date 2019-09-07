Life as a Servie (LasS)

LaaS allows you to create sessions of Conway's Game Of Life
(also known simply as Life) and watch them progress over time.
You can also browse through other user's sessions and observe them.
In case you don't know what Game Of Life is, here's a
wikipedia link - https://en.wikipedia.org/wiki/Conway's_Game_of_Life

The repo contains both the server and the client for the application.
As you might have guessed the server is used to host sessions of the game
and you can create and access them through the client.

Usage

* Client
    The client is simple: you describe a request at a prompt and it is sent to
    the server. If there is a problem with yout input an error message will be
    displayed.

    - `connect` to server
      args: server_name@ip

    - `disconnect` from server
      args: none

    - `start` new session
      args: session name, [predefined config name, path_to_config]
      See below for details on the configuratins.

    - `kill` session
      args: none
      Permanently removes a session from the server.

    - `stop` session
      command: pause
      args: session name

    - `resume` session
      args: session name

    - `watch` session
      args: session name
      Continuously displays the state of the game associated with the session.
      Use Ctrl-C to stop it. This will not stop the entire client.

* Config files
      Config files are simple text files starting with a line describing
      the dimensions of the game board and an almost graphical description
      of the board itself using '*'s for live tiles and spaces for dead tiles.
      Refer to the files containing the predefined configurations for examples.
      The client supports several predefined configurations (list of their
      names give below) which can be provided as arguments to the start
      command.


* List of predefined configurations:
    - Beehive
    - Blinker
    - Pulsar
    - Beacon
    - Pulsar
    - Glider
    - MWSS

