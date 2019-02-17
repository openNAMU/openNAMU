from .tool.func import *

def acl_list_2(conn):
    curs = conn.cursor()
    
    div =   '''
        <table id="main_table_set">
            <tbody>
                <tr>
                    <td id="main_table_width_quarter">''' + load_lang('document_name') + '''</td>
                    <td id="main_table_width_quarter">''' + load_lang('document') + ''' acl</td>
                    <td id="main_table_width_quarter">''' + load_lang('discussion') + ''' acl</td>
                    <td id="main_table_width_quarter">''' + load_lang('acl_required') + '''</td>
    '''
    
    curs.execute("select title, dec, dis, view, why from acl where dec = 'admin' or dec = 'user' or dis = 'admin' or dis = 'user' or view = 'admin' or view = 'user' order by title desc")
    list_data = curs.fetchall()
    for data in list_data:
        if not re.search('^user:', data[0]) and not re.search('^file:', data[0]):
            acl = []
            for i in range(1, 4):
                if data[i] == 'admin':
                    acl += [load_lang('admin')]
                else:
                    acl += [load_lang('member')]

            div +=  '''
                <tr>
                    <td>
                        <a href="/w/''' + url_pas(data[0]) + '">' + data[0] + '''</a>
                    </td>
                    <td>''' + acl[0] + '''</td>
                    <td>''' + acl[1] + '''</td>
                    <td>''' + acl[2] + '''</td>
                </tr>
            '''
        
    div +=  '''
            </tbody>
        </table>
    '''
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('acl_document_list'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('return')]]
    ))