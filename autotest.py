import docker
import time
import subprocess
import re

def main():
  for i in range(0,200):
    print("Launching attempt {}".format(i))
    launch()
    print("Sleeping for 15 seconds...")
    time.sleep(15)
    client = docker.from_env()
    container = client.containers.get('quicksilver_quicksilver_1')
    for logitem in container.logs(stream=True):
      restart = process_log(logitem)
      if restart:
        break


def launch():
   subprocess.Popen(
        ["make", "test-docker"],
        close_fds=True, stdin=subprocess.DEVNULL, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

def process_log(input):
  input = re.sub(r"\x1B\[([0-9]{1,3}(;[0-9]{1,2})?)?[mGK]", "", input.decode('utf-8'))
  #print("Processing log: '{}'".format(input).rstrip())
  heightMatch = re.search("indexed block height=([0-9]*) module=txindex", input, flags=0)
  apphashMatch = re.search("wrong Block.Header.AppHash", input, flags=0)
  if apphashMatch != None:
    print("Matched apphash error")
    return sys.exit(1)

  if heightMatch != None:
    if int(heightMatch.group(1))%10 == 0:
      print("Matched height {}".format(heightMatch.group(1)))
    if int(heightMatch.group(1)) >= 200:
      return True

  return False

main()
