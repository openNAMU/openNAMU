from .tool.func import *

def bbs_w_comment_tool(bbs_num = '', post_num = '', comment_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        data = ''
        
        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)
        
        data += '''
            <h2>''' + load_lang('tool') + '''</h2>
            <ul class="opennamu_ul">
                <li><a href="/bbs/raw/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '/' + url_pas(comment_num) + '">' + load_lang('raw') + '''</a></li>
                <li><a href="/bbs/edit/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '/' + url_pas(comment_num) + '">' + load_lang('edit') + '''</a></li>
            </ul>
        '''

        if admin_check() == 1:
            data += '''
                <h3>''' + load_lang('owner') + '''</h2>
                <ul class="opennamu_ul">
                    <li><a href="/bbs/delete/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '/' + url_pas(comment_num) + '">' + load_lang('delete') + '''</a></li>
                </ul>
            '''

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('bbs_comment_tool'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = data,
            menu = [['bbs/w/' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '#' + url_pas(comment_num), load_lang('return')]]
        ))