import pymysql
import json
import re

json_data = open('set.json').read()
set_data = json.loads(json_data)

conn = pymysql.connect(user = set_data['user'], password = set_data['pw'], charset = 'utf8mb4', db = set_data['db'])
curs = conn.cursor(pymysql.cursors.DictCursor)   

db_pas = pymysql.escape_string

r_ver = '2.1.5'

curs.execute('select data from other where name = "version"')
version = curs.fetchall()
if(version):
    t_ver = re.sub('\.', '', version[0]['data'])
    t_ver = re.sub('[a-z]$', '', t_ver)
    r_t_ver = re.sub('\.', '', r_ver)
    r_t_ver = re.sub('[a-z]$', '', r_t_ver)
    if(int(t_ver) <= int(r_t_ver)):
        curs.execute("update other set data = '" + db_pas(r_ver) + "' where name = 'version'")
else:
    curs.execute("insert into other (name, data) value ('version', '" + db_pas(r_ver) + "')")
    t_ver = 0
    
curs.execute('select name from alist limit 1')
getalist = curs.fetchall()
if(getalist and int(t_ver) < 204):
    curs.execute("delete from alist where name = 'owner'")
    curs.execute("delete from alist where name = 'admin'")

if(int(t_ver) < 202 or not getalist):
    curs.execute("insert into alist (name, acl) value ('owner', 'owner')")
    curs.execute("insert into alist (name, acl) value ('admin', 'ban')")
    curs.execute("insert into alist (name, acl) value ('admin', 'mdel')")
    curs.execute("insert into alist (name, acl) value ('admin', 'toron')")
    curs.execute("insert into alist (name, acl) value ('admin', 'check')")
    curs.execute("insert into alist (name, acl) value ('admin', 'acl')")
    
if(int(t_ver) < 203):
    curs.execute('select title from topic limit 1')
    top_yes = curs.fetchall()
    if(top_yes):
        curs.execute('rename table topic to old_topic')
        curs.execute('rename table distop to old_distop')
        
        curs.execute('create table topic(id text, title text, sub text, data longtext, date text, ip text, block text, top text)')
        
        curs.execute('select * from old_topic')
        topic_old = curs.fetchall()
        if(topic_old):
            i = 0
            for move_topic in topic_old:
                curs.execute("select id from distop where id = '" + db_pas(move_topic['id']) + "' and title = '" + db_pas(move_topic['title']) + "' and sub = '" + db_pas(move_topic['sub']) + "'")
                distop = curs.fetchall()
                if(distop):
                    top = 'O'
                else:
                    top = ''
                    
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) value ('" + db_pas(move_topic['id']) + "', '" + db_pas(move_topic['title']) + "', '" + db_pas(move_topic['sub']) + "', '" + db_pas(move_topic['data']) + "', '" + db_pas(move_topic['date']) + "', '" + db_pas(move_topic['ip']) + "', '" + db_pas(move_topic['block']) + "', '" + db_pas(top) + "')")
    
conn.commit()
conn.close()

print('종료 하려면 아무 키나 누르시오.')
a = input()