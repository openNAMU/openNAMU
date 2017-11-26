import re
from urllib import parse

def url_pas(data):
    return(parse.quote(data).replace('/','%2F'))

def redirect_pas(data, title, backlink):    
    d_re = re.findall('\r\n#(?:redirect|넘겨주기) ((?:(?!\r|\n|%0D).)+)', data)
    for d in d_re:
        view = d.replace('\\', '')    
            
        sh = ''
        s_d = re.search('#((?:(?!x27;|#).)+)$', d)
        if(s_d):
            href = re.sub('#((?:(?!x27;|#).)+)$', '', d)
            sh = '#' + s_d.groups()[0]
        else:
            href = d
            
        a = href.replace('&#x27;', "'").replace('&quot;', '"').replace('\\\\', '<slash>').replace('\\', '').replace('<slash>', '\\')
        backlink += [[title, a, 'redirect']]
        data = re.sub('\r\n#(?:redirect|넘겨주기) ((?:(?!\r|\n|%0D).)+)', '<meta http-equiv="refresh" content="0;url=/w/' + url_pas(a) + '/from/' + url_pas(title) + sh + '" />', data, 1)
        
    return([data, backlink])