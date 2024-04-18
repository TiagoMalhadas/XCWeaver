import time
from tqdm import tqdm
import requests
import random
import string


def metrics():
  from plumbum.cmd import xcweaver
  import re

  pattern = re.compile(r'^.*│.*│.*│.*│\s*(\d+\.?\d*)\s*│.*$', re.MULTILINE)

  def get_filter_metrics(metric_name):
    return xcweaver['multi', 'metrics', metric_name]()

  # wkr2 api
  inconsitencies_metrics = get_filter_metrics('sn_inconsistencies')
  inconsistencies_count = sum(int(value) for value in pattern.findall(inconsitencies_metrics))

  print(f"number of inconsistencies: {inconsistencies_count}")


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
    print(f"number of requests: {n_requests}")

    #TO-DO
    #Get metrics

if __name__ == "__main__":
    run_test(30)
    metrics()
    print(f"[INFO] done!")
    
