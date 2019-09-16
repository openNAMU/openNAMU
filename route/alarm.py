from .tool.func import *

def alarm_2():
    

    data = '<ul>'    
    
    sqlQuery("select data, date from alarm where name = ? order by date desc", [ip_check()])
    data_list = sqlQuery("fetchall")
    if data_list:
        data = '<a href="/del_alarm">(' + load_lang('delete') + ')</a><hr class=\"main_hr\">' + data

        for data_one in data_list:
            data += '<li>' + data_one[0] + ' (' + data_one[1] + ')</li>'
    
    data += '</ul>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('notice'), wiki_set(), custom(), other2([0, 0])],
        data = data,
        menu = [['user', load_lang('return')]]
    ))