from typing import Dict


class ThreadData:
    

    
    def __init__(self) -> None:
        pass
    
    def __init__(self, data) -> None:
        self.uuid = data["uuid"]
        self.date = data["date"]
        self.file = data["file"]
        self.responses = data["responses"]
        self.images = data["images"]
        self.teaser = data["teaser"]
        self.images_url = data["image-url"]
        self.teaser_thumb_width = data["teaser-thumb-width"]
        self.teaser_thumb_height = data["teaser-thumb-height"]
        self.board_name = data["board-name"]
        

    def toSqliteInsert(self):
        return f'INSERT OR IGNORE INTO `threads` (`id`, `date`, `file`, `responses`, `images`, `teaser`, `images-url`, `board`) VALUES ({self.uuid}, {self.date}, "{self.file}", {self.responses}, {self.images}, "{self.teaser}", "{self.images_url}", "{self.board_name}");'
    