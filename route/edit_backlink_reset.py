from .tool.func import *

async def edit_backlink_reset(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change("select data from data where title = ?"), [name])
        old = curs.fetchall()
        if old:
            render_set(conn, 
                doc_name = name,
                doc_data = old[0][0],
                data_type = 'backlink'
            )

        return redirect(conn, '/xref/' + url_pas(name))