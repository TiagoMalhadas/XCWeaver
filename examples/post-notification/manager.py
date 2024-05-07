#!/usr/bin/env python3
import argparse
import datetime
import sys
import time
from plumbum import FG
from tqdm import tqdm
import requests
import random
import string


def metrics():
  from plumbum.cmd import xcweaver
  import re

  timestamp = datetime.datetime.now().strftime("%Y-%m-%d_%H:%M:%S")

  pattern = re.compile(r'^.*│.*│.*│.*│\s*(\d+\.?\d*)\s*│.*$', re.MULTILINE)

  def get_filter_metrics(metric_name):
    return xcweaver['multi', 'metrics', metric_name]()

  # wkr2 api
  inconsitencies_metrics = get_filter_metrics('sn_inconsistencies')
  inconsistencies_count = sum(int(value) for value in pattern.findall(inconsitencies_metrics))
  requests_metrics = get_filter_metrics('requests')
  requests = sum(int(value) for value in pattern.findall(requests_metrics))
  pc_inconsistencies = "{:.2f}".format((inconsistencies_count / requests) * 100)
  post_notification_duration_metrics = get_filter_metrics('sn_post_notification_duration_ms')
  post_notification_duration_metrics_values = pattern.findall(post_notification_duration_metrics)
  post_notification_duration_avg_ms = sum(float(value) for value in post_notification_duration_metrics_values if value != 0)/2 if post_notification_duration_metrics_values else 0
  write_post_duration_metrics = get_filter_metrics('sn_write_post_duration_ms')
  write_post_duration_metrics_values = pattern.findall(write_post_duration_metrics)
  print(write_post_duration_metrics_values)
  write_post_duration_avg_ms = sum(float(value) for value in write_post_duration_metrics_values if value != 0)/2 if write_post_duration_metrics_values else 0
  notifications_sent_metrics = get_filter_metrics('sn_notificationsSent')
  notifications_sent = sum(int(value) for value in pattern.findall(notifications_sent_metrics))
  notifications_received_metrics = get_filter_metrics('notificationsReceived')
  notifications_received = sum(int(value) for value in pattern.findall(notifications_received_metrics))
  percentage_notifications_received = "{:.2f}".format((notifications_received / notifications_sent) * 100)
  read_post_duration_metrics = get_filter_metrics('sn_read_post_duration_ms')
  read_post_duration_metrics_values = pattern.findall(read_post_duration_metrics)
  print(read_post_duration_metrics_values)
  read_post_duration_avg_ms = sum(float(value) for value in read_post_duration_metrics_values if value != 0)/2 if read_post_duration_metrics_values else 0
  queue_duration_metrics = get_filter_metrics('sn_queue_duration_ms')
  queue_duration_metrics_values = pattern.findall(queue_duration_metrics)
  print(queue_duration_metrics_values)
  queue_duration_avg_ms = sum(float(value) for value in queue_duration_metrics_values if value != 0)/2 if queue_duration_metrics_values else 0

  results = f"""
    # requests:\t\t\t{requests}
    # received notifications @ US:\t{notifications_received} ({percentage_notifications_received}%)
    # inconsistencies @ US:\t\t{inconsistencies_count}
    % inconsistencies @ US:\t\t{pc_inconsistencies}%
    > avg. post notification duration:\t{post_notification_duration_avg_ms}ms
    > avg. write post duration:\t\t{write_post_duration_avg_ms}ms
    > avg. read post duration:\t\t{read_post_duration_avg_ms}ms
    > avg. queue duration @ US:\t\t{queue_duration_avg_ms}ms
  """
  print(results)

  # save file if we ran workload
  if timestamp:
    filepath = f"evaluation/local/{timestamp}_metrics.txt"
    with open(filepath, "w") as f:
      f.write(results)
    print(f"[INFO] evaluation results saved at {filepath}")


def run_test(duration):
    import threading
    def tqdm_progress(duration):
        print(f"[INFO] running workload for {duration} seconds...")
        for _ in tqdm(range(int(duration))):
            time.sleep(1)

    progress_thread = threading.Thread(target=tqdm_progress, args=(duration,))
    progress_thread.start()

    start_time = time.time()
    n_requests = 0
    url = "http://localhost:12345/post_notification"

    #execute requests for the given time
    while True:
        n_requests += 1
        post = ''.join(random.choice(string.ascii_letters + string.digits) for _ in range(15))
        params = {
            "post": post
        }

        response = requests.get(url, params=params)

        if response.status_code != 200:
            print(f"Request failed with status code {response.status_code}.")

        elapsed_time = time.time() - start_time
        if elapsed_time >= duration:
            break

        time.sleep(0.01)

    progress_thread.join()
    progress_thread = threading.Thread(target=tqdm_progress, args=(180,))
    progress_thread.start()
    progress_thread.join()

def local_wrk2():
  run_test(30)
  metrics()

def local_metrics():
  metrics()

def local_storage_run():
  from plumbum.cmd import docker_compose, docker
  docker_compose['up', '-d'] & FG
  print("[INFO] waiting 30 seconds for storages to be ready...")
  for _ in tqdm(range(30)):
      time.sleep(1)
  docker['exec', '-it', 'mongo-1', 'mongo', '--eval', "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-1:27017'}, {_id: 1, host: 'mongo-2:27017'}]})"] & FG
  print("[INFO] getting replica set status")
  docker['exec', '-it', 'mongo-1', 'mongo', '--eval', "rs.status()"] & FG

def local_storage_info():
  print("[INFO] nothing to be done for local")
  exit(0)

def local_storage_clean():
  from plumbum.cmd import docker_compose
  docker_compose['down'] & FG

if __name__ == "__main__":
  main_parser = argparse.ArgumentParser()
  command_parser = main_parser.add_subparsers(help='commands', dest='command')

  commands = [
    # datastores
    'storage-run', 'storage-info', 'storage-clean',
    # eval
    'wrk2', 'metrics',
  ]

  for cmd in commands:
    parser = command_parser.add_parser(cmd)
    parser.add_argument('--local', action='store_true', help="Running in localhost")
    parser.add_argument('--gcp', action='store_true',   help="Running in gcp")

  args = vars(main_parser.parse_args())
  command = args.pop('command').replace('-', '_')

  local = args.pop('local')
  gcp = args.pop('gcp')

  if local and gcp or not local and not gcp:
    print("[ERROR] one of --local or --gcp flgs needs to be provided")
    exit(-1)

  if local:
    command = 'local_' + command
  elif gcp:
    #load_gcp_profile()
    command = 'gcp_' + command

  print(f"[INFO] ----- {command.upper().replace('_', ' ')} -----\n")
  getattr(sys.modules[__name__], command)(**args)

  print(f"[INFO] done!")
    
