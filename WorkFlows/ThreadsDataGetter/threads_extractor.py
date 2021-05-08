
from prefect.schedules import IntervalSchedule
from bs4 import BeautifulSoup as soup
from prefect import task, Flow
from typing import Dict, List
from os import getcwd, getenv
from os.path import exists
from ThreadClass import ThreadData
import requests, json, datetime, string
import re, sqlite3

CHAN_NAME = ["b", "aco", "d", "h", "s", "hc"]
SAVE_FILENAME = 'threads_data.json'
DOWNLOAD_SERVER = "http://0.0.0.0:5006" if not getenv("DOWNLOAD_SERVER") else getenv("DOWNLOAD_SERVER")


sqlite_connection = sqlite3.connect('threads_data.db')

valid_characters = set(string.ascii_letters + string.digits + " ")

print(f"CURRENT DIR: {getcwd()}")
print(f"{DOWNLOAD_SERVER=}")

def alertFailure(obj, old_state, new_state):
    if new_state.is_failed():
        print(f"FAILURE: {obj=}\n{old_state=}\n{new_state}")

@task(name="download relevant threads")
def downloadRelevantThread(board_table:Dict, board_name:str):
    for board_thread in board_table:
        if any(map(lambda x: x in board_thread['teaser'], MEANINGFUL_NAMES)):
            sendDownloadRequest(board_thread, board_name)
        


def get4ChanUrl(name:str) -> str:
    chan_url = f"https://boards.4chan.org/{name}/catalog"
    
    #cheking if board exists
    cursor = sqlite_connection.cursor()
    cursor.execute(f"SELECT `uuid` FROM `boards` WHERE `uuid`='{name}';")
    if cursor.fetchone() == None: 
        #board is not registred
        cursor.execute(f"INSERT OR REPLACE INTO `boards`(`uuid`, `url`) VALUES ('{name}', '{chan_url}');")
        sqlite_connection.commit()
        cursor.close()
    
    return chan_url

@task(state_handlers=[alertFailure])
def getChan(name) -> str:
    response = requests.get(get4ChanUrl(name))
    if response.ok:
        return response.text
    raise requests.ConnectionError

@task
def parseTable(content: str) -> Dict:
    data_soup = soup(content,'html.parser')
    threads_tag = data_soup.findAll('script')[4]
    table = "{" + re.search(r"var\scatalog\s=\s\{.+\};", threads_tag.string)[0].split("{",1)[1]
    json_data = json.loads(table[:-1])
    return json_data

def parseThread(thread_url:str) -> List:
    response = requests.get(thread_url)
    links = []
    if response.ok:
        data_soup = soup(response.text, 'html5')
        image_tags = data_soup.select(".fileText a")
        for it in image_tags:
            if it.has_attr("href"):
                links.append(f'https:{it["href"]}')
    return links

def transformTeaser(teaser: str, word_max_count:int=5, word_max_length:int=20) -> str:
    teaser = "".join(filter(lambda x: x in valid_characters, teaser))
    new_teaser = []
    for w in teaser.split(" "):
        if len(new_teaser) > 5:
            break
        if len(w) > 20:
            continue
        new_teaser.append(w)
    return " ".join(new_teaser)

def sendDownloadRequest(thread_data:Dict, board_name):
    thread_url = f"https://boards.4chan.org/{board_name}/thread/{thread_data['uuid']}/"
    images = parseThread(thread_url)
    request_form = {
                    "name": transformTeaser(thread_data['teaser']),
                    "path": ".",
                    "files":json.dumps(images)
                    }
    response = requests.post(f"{DOWNLOAD_SERVER}/downloads", request_form)
    if response.ok:
        print(f"send download request for: {thread_url}")
    else:
        print(f"something went wrong with download '{thread_data}'")
    




@task
def extract(json_table: Dict, board_name: str) -> Dict:
    meaningful_table = [] 
    for thread_uuid, chan_thread in json_table["threads"].items():
        if "imgurl" in chan_thread:
            table_row = ThreadData({                
                "uuid": int(thread_uuid),
                "date": chan_thread["date"],
                "file": chan_thread["file"],
                "responses": chan_thread["r"],
                "images": chan_thread["i"],
                "teaser": transformTeaser(chan_thread["teaser"].lower(), word_max_count=50),
                "image-url": chan_thread["imgurl"],
                "teaser-thumb-width": chan_thread["tn_w"],
                "teaser-thumb-height": chan_thread["tn_h"],
                "board-name": board_name
            }
            )
            meaningful_table.append(table_row)
    return meaningful_table



@task
def saveThreadsToSqlite(board_tabel: List[ThreadData]) -> None:
    sqlite_cursor = sqlite_connection.cursor()
    for board_thread in board_tabel:
        sqlite_cursor.execute(board_thread.toSqliteInsert())
    sqlite_connection.commit()
    sqlite_cursor.close()

schedule = IntervalSchedule(interval=datetime.timedelta(minutes=30))

with Flow("get-chan-threads") as the_flow:

    for chan_name in CHAN_NAME:
        c = getChan(chan_name)
        r = parseTable(c)
        r= extract(r, chan_name)
        
        
        
        # download service down here
        # downloadRelevantThread(bt, chan_name)
        saveThreadsToSqlite(r)
the_flow.run()
sqlite_connection.close()
