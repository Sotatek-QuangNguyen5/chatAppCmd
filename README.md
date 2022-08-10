# chatAppCmd

run server :\
cd server\
go run server.go\


run client :\
cd client\
go run client.go\

+) tao nhieu client\
+) moi khi chay file client.go\
+) dau tien se nhap ten (neu de trong se phai nhap lai), sau do se gui connect den server, server se tao ra user_id duy nhat cho user do va luu user do lai
khi ms chay user se ko trong 1 cuoc tro chuyen nao, do do khi chat, client se bao la ban dang ko o trong room nao, vui long nhap command
khi nhap command se hien ra dong hoi co phai ban dang nhap lenh ko, neu tra loi la yes, no se hieu la 1 lenh, neu ko no se hieu la 1 tin nhan
+) co 8 loai command\
-) command help : hien tat ca cac command\
-) command exit : thoai ra khoi cuoc tro chuyen hien tai\
-) command join : nhap vao them roomID, join vao 1 cuoc tro chuyen bang roomID, neu no ko ton tai, server se tra ve cho client thong bao room ko ton tai, neu room do la -) private, thong bao lai cho client ko join vao dc, private room la chat 2 nguoi vs nhau, neu roomID la room hien tai, thong bao dang o trong phong roi, khi join vao 1 phong khac phong hien tai, ban se thoat ra khoi phong hien tai\
-) command create : nhap them roomID, tao phong chat public, nhieu nguoi join vao duoc, neu nhap roomID da ton tai, server se bao lai cho client room da ton tai, create xong ban se join vao phong do luon va thoat ra khoi phong hien tai\
-) command list_user : show tat ca cac user dang hoat dong, ko tinh \
-) command list_room : show tat ca cac room dang hoat dong, ko tinh room minh dang o trong \
-) command chat_with : nhap them user_id, neu nguoi do dang ko chat vs ai, thi se tao ra 1 phong private, chi co 2 nguoi chat vs nhau, nguoi khac khong join vao dc, neu user_id ko ton tai se bao ve cho client, neu user do dang chat trong phong khac se thong bao la nguoi do dang o phong khac roi \

