# bindu-beta
Simple api development go framework top on gin
bindu new  --app Blank --port 8080 #true

bindu generate Help #true

bindu generate Scaffold User Email:String UserName:String Phone:String --hasOne Profile 
bindu generate Model Profile FirstName:String LastName:String Image:String  
bindu generate Controller Profile FirstName:String LastName:String Image:String  

bindu generate Scaffold Item Location:String Type:String Image:String Title:String Date:String 

bindu generate Model User --hasMany Item #false

bindu create User Username:Sohag Password:11111 Role:Admin 
bindu generate Routes User --middleware auth --group v1 
bindu generate Routes Item --middleware auth --group v1 
bindu generate Controller User

bindu create Policy Alice:P Sub:Admin Obj:/Api/* Act:* 

bindu generate Scaffold Blog Title:String Description:String --hasMany Comment --belongsTo User 

bindu generate Scaffold Comment Subject:String Body:String 


bindu generate Scaffold Contact FirstName:String LastName:String Subject:String Description:String

bindu generate Model Item Description:String 

bindu generate Model Item Mobile:String 

bindu db Migrate


bindu create Policy Alice:P Sub:Subscriber Obj:/Api/V1/* Act:* 

bindu create User Username:Sohag Password:11111 Role:Admin 

