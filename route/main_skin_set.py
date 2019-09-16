from .tool.func import *

def main_skin_set_2():
    
    
    data = flask.make_response(re_error('/error/5'))

    sqlQuery("select data from other where name = 'language'")
    main_data = sqlQuery("fetchall")

    data.set_cookie('language', main_data[0][0])

    sqlQuery('select data from user_set where name = "lang" and id = ?', [ip_check()])
    user_data = sqlQuery("fetchall")
    if user_data:
        data.set_cookie('user_language', user_data[0][0])
    else:
        data.set_cookie('user_language', main_data[0][0])

    return data