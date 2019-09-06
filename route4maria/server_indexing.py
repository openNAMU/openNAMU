from .tool.func import *
import pymysql

def server_indexing_2(conn):
    curs = conn.cursor()

    if admin_check() != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        admin_check(None, 'indexing')

        curs.execute("select name from sqlite_master where type = 'index'")
        data = curs.fetchall()
        if data:
            for delete_index in data:
                print('Delete : ' + delete_index[0])

                sql = 'drop index if exists ' + delete_index[0]
                
                try:
                    curs.execute(sql)
                except:
                    pass
        else:
            curs.execute("select name from sqlite_master where type in ('table', 'view') and name not like 'sqlite_%' union all select name from sqlite_temp_master where type in ('table', 'view') order by 1;")
            for table in curs.fetchall():            
                curs.execute('select sql from sqlite_master where name = %s', [table[0]])
                cul = curs.fetchall()
                
                r_cul = re.findall('(?:([^ (]*) text)', str(cul[0]))
                
                for n_cul in r_cul:
                    print('Create : index_' + table[0] + '_' + n_cul)

                    sql = 'create index index_' + table[0] + '_' + n_cul + ' on ' + table[0] + '(' + n_cul + ')'
                    try:
                        curs.execute(sql)
                    except:
                        pass

        
        
        return redirect()  
    else:
        curs.execute("select name from sqlite_master where type = 'index'")
        data = curs.fetchall()
        if data:
            b_data = load_lang('delete')
        else:
            b_data = load_lang('create')

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('indexing'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <button type="submit">''' + b_data + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('return')]]
        ))   