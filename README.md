# Quicksilver

The Cosmos Liquid Staking Zone

## Software Dependencies

1. The Go programming language - https://go.dev/
2. Git distributed version control - https://git-scm.com/
3. Docker - https://www.docker.com/get-started/
4. GNU Make - https://www.gnu.org/software/make/

Make sure that the above software is installed on your system. Follow the instructions for your particular platform or use your preferred platform package manager;

In addition install `jq` (a command line JSON processor):

 - Debian based systems:  
`apt-get install jq`

 - Arch based systems:  
`pacman -S jq`

 - Mac based systems:  
`brew install jq`

 - Windows based systems (using [Chocolatey NuGet](https://chocolatey.org/)):  
`chocolatey install jq`

## Clone & Run Quicksilver (dev)

_NB!! Use a fork of the repository when you plan to create Pull Requests;_

Clone the repository from GitHub and enter the directory:

    git clone https://github.com/ingenuity-build/quicksilver.git
    cd quicksilver

Then run:

    make build-docker
    make test-docker

For subsequent tests run the following if you want to start with fresh state:

    make build-docker
    make test-docker-regen
