from .tool.func import *

def bbs_w(bbs_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect('/bbs/main')
        
        bbs_name = db_data[0][0]
        bbs_num_str = str(bbs_num)

        data = ''
        data += '''
            <table id="main_table_set">
                <tr id="main_table_top_tr">
                    <td id="main_table_width">''' + load_lang('editor') + '''</td>
                    <td id="main_table_width">''' + load_lang('time') + '''</td>
                    <td id="main_table_width">''' + load_lang('last_comment_time') + '''</td>
                </tr>
        '''
        
        temp_id = ''
        temp_dict = {}

        curs.execute(db_change('select set_name, set_data, set_code, set_id from bbs_data where set_id like ? order by set_code + 0 desc'), [bbs_num])
        db_data = curs.fetchall()
        db_data = list(db_data) if db_data else []

        for for_a in db_data + [['', '', '']]:
            if temp_id != for_a[2]:
                if temp_id != '':
                    curs.execute(db_change('select count(*) from bbs_data where set_name = "comment_date" and (set_id = ? or set_id like ?) order by set_code + 0 desc'), [bbs_num_str + '-' + temp_dict['code'], bbs_num_str + '-' + temp_dict['code'] + '-%'])
                    db_data = curs.fetchall()
                    comment_count = str(db_data[0][0]) if db_data else '0'

                    curs.execute(db_change('select set_data from bbs_data where set_name = "comment_date" and (set_id = ? or set_id like ?) order by set_data desc limit 1'), [bbs_num_str + '-' + temp_dict['code'], bbs_num_str + '-' + temp_dict['code'] + '-%'])
                    db_data = curs.fetchall()
                    last_comment_date = db_data[0][0] if db_data else '0'

                    data += '''
                        <tr>
                            <td>''' + ip_pas(temp_dict['user_id']) + '''</td>
                            <td>''' + temp_dict['date'] + '''</td>
                            <td>''' + last_comment_date + '''</td>
                        </tr>
                        <tr>
                            <td colspan="3">
                                <a href="/bbs/w/''' + bbs_num_str + '/' + temp_dict['code'] + '">' + html.escape(temp_dict['title']) + '''</a> 
                                (''' + comment_count + ''') 
                            </td>
                        </tr>
                    '''

                temp_id = for_a[2]
                temp_dict['code'] = for_a[2]

            temp_dict[for_a[0]] = for_a[1]

        data += '</table>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [bbs_name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('bbs') + ')', 0])],
            data = data,
            menu = [['bbs/main', load_lang('return')], ['bbs/edit/' + bbs_num_str, load_lang('add')], ['bbs/set/' + bbs_num_str, load_lang('bbs_set')]]
        ))