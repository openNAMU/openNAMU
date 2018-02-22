import re

def text_help(data):    
    data = re.sub("&#x27;&#x27;&#x27;(?P<in>(?:(?!&#x27;&#x27;&#x27;).)*)&#x27;&#x27;&#x27;", '<b>\g<in></b>', data)
    data = re.sub("&#x27;&#x27;(?P<in>(?:(?!&#x27;&#x27;).)*)&#x27;&#x27;", '<i>\g<in></i>', data)
    data = re.sub('(?:~~|--)(?P<in>(?:(?!~~|--).)+)(?:~~|--)', '<s>\g<in></s>', data)
    data = re.sub('__(?P<in>.+?)__(?!_)', '<u>\g<in></u>', data)
    data = re.sub('\^\^(?P<in>.+?)\^\^(?!\^)', '<sup>\g<in></sup>', data)
    data = re.sub(',,(?P<in>.+?),,(?!,)', '<sub>\g<in></sub>', data) 
    data = re.sub('{{\|(?P<in>(?:(?:(?:(?!\|}}).)*)(?:\n?))+)\|}}', '<table><tbody><tr><td>\g<in></td></tr></tbody></table>', data)

    while 1:
        if re.search("\n-{4,9}( *)\r\n", data):
            data = re.sub("\n-{4,9}( *)\r\n", "\n<hr style='margin-top: 10px; margin-bottom: -10px;'>\n", data, 1)
        else:
            break
        
    return data