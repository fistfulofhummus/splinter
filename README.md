
# Splinter
A implant I am designing to work with netcat.
  

## Usage/Examples
Before building the implant don't forget to edit the code and insert your TCP listner address.

    git clone https://github.com/fistfulofhummus/splinter.git
    go build
    nc -nlvp port

Now just execute the implant on a Windows machine. If you want to terminate execution cleanly, type stop into your netcat terminal.

## TO-DO

 - [x]  Add a menu a la meterpreter
 - [ ]  Load and execute shellcode
 - [ ]  Upload/Download
 - [ ] Persistence
 - [ ] Keylogging
 - [ ] RickRoll

## Notes:
I will not be responsible for any illegal activity conducted with this code.
The implant uses raw TCP. OPSEC is a non-concern in this project. IDS/IPS systems should sniff it out.
