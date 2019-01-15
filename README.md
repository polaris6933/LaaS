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

Both the server and the client use a cli with a fairly simple way of
issuing commands


* Client
    The client has two modes: command mode and viewing mode.
    In command mode you can connect to a server hosting sessions of the game
    and issue commands for interacting with these sessions
    In viewing mode you watch a session of the game unfold and have a
    degree of control over it using one button shortcuts. Most of these
    actions are authorized via a password set during the creation of the
    session. For convinience you can unlock a session using its password and
    you won't have to authorize individual actions. Newly created sessions
    are automatically unlocked.
    The actual names of the commands are enclosed in `` in ths manual.

    - `connect` to server
      args: server_name@ip

    - `disconnect` from server
      args: none

    - `start` new session
      args: session name, [predefined config name, path_to_config]
      See below for details on the configuratins. If no argument is given a
      random one will be genrated. To finalize the creation you will be
      prompted for a password.

    - `kill` session
      args: none
      viewing mode shortcut: k

    - `unlock` session
      args: none
      viewing mode shortcut: u
      You will need to provide the password of the session. Afterwards
      you will have access to all its functionallity.

    - `save` session
      args: session name, path to export to
      viewing mode shortcut: s

    - `pause` session
      command: pause
      args: session name
      viewing mode shortcut: p

    - `resume` session
      args: session name
      viewing mode shortcut: r


* Server
    The server is pretty simple since the interaction is achieved via the
    client.

    - `start` server
      args: name

    - `stop` server
      args: none

    - `rename` server
      args: new name

    - `list` sessions
      args: none

    - `kill` session
      args: session name

    - `save` session
      args: session name
      It makes sense to save sessions before deleting them so they can
      be restored if needed. Deleted session will show up when using
      `list`.

    - `restore` session
      args: session name


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

