import re

def text_help(data):    
    data = re.sub("&#x27;&#x27;&#x27;(?P<in>(?:(?!&#x27;&#x27;&#x27;).)*)&#x27;&#x27;&#x27;", '<b>\g<in></b>', data)
    data = re.sub("&#x27;&#x27;(?P<in>(?:(?!&#x27;&#x27;).)*)&#x27;&#x27;", '<i>\g<in></i>', data)
    data = re.sub('(?:~~|--)(?P<in>(?:(?!~~|--).)+)(?:~~|--)', '<s>\g<in></s>', data)
    data = re.sub('__(?P<in>.+?)__(?!_)', '<u>\g<in></u>', data)
    data = re.sub('\^\^(?P<in>.+?)\^\^(?!\^)', '<sup>\g<in></sup>', data)
    data = re.sub(',,(?P<in>.+?),,(?!,)', '<sub>\g<in></sub>', data)

    if re.search('&lt;math&gt;((?:(?!&lt;\/math&gt;).)*)&lt;\/math&gt;', data):
        data = re.sub('&lt;math&gt;(?P<in>(?:(?!&lt;\/math&gt;).)*)&lt;\/math&gt;', '[math]\g<in>[/math]', data)
        data += '<script type="text/x-mathjax-config">MathJax.Hub.Config({tex2jax: { inlineMath: [ [\'[math]\', \'[/math]\'] ] } });</script> \
                <script src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.3/MathJax.js?config=TeX-AMS_CHTML"></script>'
 
    data = re.sub('{{\|(?P<in>(?:(?:(?:(?!\|}}).)*)(?:\n?))+)\|}}', '<table><tbody><tr><td>\g<in></td></tr></tbody></table>', data)
    data = re.sub("-{4,11}", "<hr>", data)
        
    return data