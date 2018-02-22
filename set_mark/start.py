import re

def start(data):    
    data = data.replace('\\r', '&#92;r')
    data = data.replace('\\n', '&#92;n')
    data = re.sub("\n", "\r\n", re.sub("\r\n", "\n", data))
    data = '\r\n' + data + '\r\n'
    data = re.sub('\n \|\|', '\n||', data)
        
    return data