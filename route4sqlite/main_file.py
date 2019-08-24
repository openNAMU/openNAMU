from .tool.func import *

def main_file_2(conn, data):
    curs = conn.cursor()

    if re.search('\.txt$', data):
        return flask.send_from_directory('./', data)
    else:
        return redirect('/w/' + url_pas(wiki_set(2)))