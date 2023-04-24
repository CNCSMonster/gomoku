
let buttons=[]
let chesses=[]  //array保存棋子
let player1ID=""
let gameId=-1;  //gameid是向服务器发送请求后得到的id,用-1表示游戏还没开始
let player1Password;
let =document.querySelector(".invite-button");
let serverID="43.136.17.142"
// let serverID="localhost"
const server='http://'+serverID+':6363/gomoku';
let isAbled=false;
let boardCase={
    curPlayer:1,    //用1表示自己，用2表示对面
    winner:-1,  //用-1表示没有人胜利
    chesses:[]  //用过-1，0,1,2分别表示hover时棋盘，棋盘,，本地玩家，远程玩家
};
/*

游戏流程:
玩家输入自己账号密码以及对手的账号
点击邀请按钮创建房间或者进入房间(如果房间存在就是进入房间，如果房间不存在就是创建房间)
当房间准备好，(两个选手均到场时,开始游戏)
（开始账号不存在的时候点击invite会创建账号)
先邀请的人先进入房间，执黑先行。
后邀请的人后进入房间，执白后行
在定时任务中刷新询问远程数据以及刷新界面和判断比赛是否结束
如果要退出房间重开,则点击exit按钮离开房间，同时远程会删除该游戏数据
如果要重新开始游戏，点击reset按钮，会向对方发送一个重开申请，如果双方都
统一重开，则游戏重开
*/


const backgroundImg="<img src=\"res/background.gif\" style=\"width:100%\">";
const whiteImg="<img src=\"res/whiteStone.gif\" style=\"width:100%\">";
const blackImg="<img src=\"res/blackStone.gif\" style=\"width:100%\">";
const backgroundDeepImg="<img src=\"res/background-deep.png\" style=\"width:100%\">";

const arr=[];

let loadImgToStorage=function(imgPath){
    let imgObj=new Image()
    imgObj.src=imgPath;
    arr.push(imgObj)
};
loadImgToStorage("res/background.gif");
loadImgToStorage("res/blackStone.gif");
loadImgToStorage("res/whiteStone.gif");
loadImgToStorage("res/background-deep.png");

loadBoardCase();
storeBoardCase();

document.addEventListener("DOMContentLoaded",()=>{   
    // 增加按钮,
    let board=document.querySelector(".chess-board");
    
    for(let i=0;i<100;i++){
        let button=document.createElement("button");
        button.innerHTML=backgroundImg;
        board.appendChild(button);
        buttons.push(button)
        const x=Math.floor(i/10);
        const y=i-10*x;
        button.addEventListener('click',()=>{
            // if(boardCase.curPlayer!==1) return;
            play(x,y);
        });
        button.addEventListener("mouseenter",()=>{
            if(!isAbled) return;
            if(boardCase.curPlayer!==1) return;
            setTimeout(() => {
                if(boardCase.chesses[x][y]===0){
                    button.innerHTML=backgroundDeepImg;
                }
            }, 50);
        });
        button.addEventListener("mouseleave",()=>{
            if(!isAbled) return;
            setTimeout(() => {
                if(boardCase.chesses[x][y]===0){
                    button.innerHTML=backgroundImg;
                }
            }, 50);
        });
    }
    // 给invite按钮增加方法
    setInterval(() => {
        if(gameId<0) return;
        fetchBoardCase();
    }, 500);
    inviteButton=document.querySelector(".invite-button");
    inviteButton.addEventListener('click',()=>{
        // 根据 invite-button中的属性来判断使用什么方法
        if(inviteButton.value==="invite"){
            invite()
        }else{
            exit()
        }

    });
    let resetButton=document.querySelector(".reset-button");
    resetButton.addEventListener("click",resetGame);

});

function resetGame(){
    if(gameId<0) return;
    //TODO 重新开始游戏，向对方发出申请,通过,
    // 发送申请给远程,这次使用参数,来反馈
    // 1. 创建XMLHttpRequest对象
    const xhr = new XMLHttpRequest();
    // 2. 定义请求的参数
    const params = 'gameID='+gameId; // 请求的参数，格式为key=value&key=value
    // 3. 因为init请求具有幂等性，所以使用PUT
    xhr.open('PUT', server +"/game"+ '?' + params, true);
    // 4. 忽略请求头部设置
    xhr.onload=function() {
        if(xhr.status===200){
            initBoardCase();
            refreshUI()
            isAbled=false;
        }else{
            alert("reset fail!")
        }
    }
    // 5. 发送请求
    xhr.send();
}

// invite 按钮需要的invite方法
function invite(){
    // 浏览器申请公钥
    // 服务端返回房间对应的公钥，
    // 浏览器通过公钥加密自己的(账号，密码，敌人，自己的公钥)发送给服务端
    // 服务端验证通过后，使用客户端的发送的公钥加密（分配的游戏id，新的公钥给客人)

    player1Account=document.getElementById("player-account").value;
    player1Password=document.getElementById("player-password").value;
    enemyAccount=document.getElementById("enemy-account").value;
    if (player1Password===""||player1Account===""||enemyAccount==="") {
        return
    }
    // 把它传递给远程
    var xhr=new XMLHttpRequest()
    xhr.open('GET', server+"/invite/"+player1Account+"/"+enemyAccount+"/"+player1Password); // 发送GET请求到服务器，请求数据的API为/api/data
    xhr.onload = function () {
        if (xhr.status === 200) {
            let responseData = xhr.responseText; // 解析服务器返回的json字符串
            // 分解得到player1ID和player2ID
            let arr= responseData.split("-");
            gameId=Number(arr[0]);
            player1ID=arr[1];
            document.querySelector(".invite-button").value="exit";
            isAbled=true;
            fetchBoardCase();
        } else {
            console.error('请求失败，状态码为：' + xhr.status);
        }
    };
    xhr.send();

}

// invite 按钮变成exit按钮后的exit方法
function exit(){
    // TODO 向远程发送信息,释放该房间的资源
    if(gameId<0) return;
    var xhr=new XMLHttpRequest()
    xhr.open('DELETE', server+"/game"+"/"+gameId+"/"+player1ID); // 发送GET请求到服务器，请求数据的API为/api/data
    xhr.onload = function () {
        if (xhr.status === 200) {
            location.href=server+"/board.html"
            inviteButton.value="invite";
            isAbled=false;
            refreshUI();
        } else {
            console.error('请求失败，状态码为：' + xhr.status);
        }
    };
    xhr.send();
}


// 下棋方法，每次下棋传入自己的坐标，如果下棋成功,返回true,如果下棋失败,返回false
function play(x,y){
    // 传入纵横坐标
    if(boardCase.winner>0) return;
    if(boardCase.chesses[x][y]!==0&&boardCase.chesses[x][y]!==-1){
        return;
    }
    // 如果游戏还没开始,不执行操作
    if(gameId<0){
        // boardCase.chesses[x][y]=boardCase.curPlayer;
        // boardCase.curPlayer=boardCase.curPlayer===1?2:1;
        // storeBoardCase();
        return;
    }
    if(!isAbled){
        setTimeout(play(x,y),100);
        return;
    }
    // 创建下棋位置对象
    const pos={
        X:x,
        Y:y
    };
    // 把下棋位置对象发送给远程
    // 创建XMLHttpRequest对象
    var xhr = new XMLHttpRequest();
    // 配置请求参数
    xhr.open('POST', server+"/game/"+gameId+"/"+player1ID, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    // 定义回调函数
    xhr.onload = function() {
        if (xhr.status == 200) {
            boardCase.chesses[x][y]=boardCase.curPlayer;
            boardCase.curPlayer=boardCase.curPlayer===1?2:1;
            storeBoardCase();
            refreshUI();
        }
        isAbled=true;
    };
    isAbled=false;
    xhr.send(JSON.stringify(pos));
}

function fetchBoardCase(){
    const xhr = new XMLHttpRequest();
    xhr.open('GET', server + "/game/" + gameId + "/" + player1ID); // 发送GET请求到服务器，请求数据的API为/api/data
    xhr.onload = function () {
        if (xhr.status === 200) {
            let newBoardCase = JSON.parse(xhr.responseText); // 解析服务器返回的json字符串
            // 判断旧的boardCase相对新的有没有改变,时间复杂度10*10
            let ifChanged=false;
            for(let i=0;i<boardCase.chesses.length&&!ifChanged;i++){
                for(let j=0;j<boardCase.chesses[i].length;j++){
                    let newChess=newBoardCase.chesses[i][j];
                    let old=boardCase.chesses[i][j];
                    if(old===newChess) continue;
                    else{
                        ifChanged=true;
                        break;
                    }
                }
            }
            if(newBoardCase.curPlayer!==boardCase.curPlayer) ifChanged=true;
            if(newBoardCase.winner!==boardCase.winner) ifChanged=true;
            if(boardCase.curPlayer===1) isAbled=true;
            if(ifChanged){
                boardCase=newBoardCase;
                refreshUI();
            }
            checkWinner();
        } else {
            console.log('请求失败，状态码为：' + xhr.status);
        }
    };
    isAbled=false;
    xhr.send();
}


function checkWinner(){
    // 检查是否有人胜利
    const outcome=document.getElementById("outcome-area");
    if ( boardCase.winner===1 || boardCase.winner===2 ){
        // 如果本地玩家胜利
        if (boardCase.winner===1&&outcome.innerText!="You Win!"){
                outcome.innerText="You Win!";
                alert("You win!");
            }
            //如果远程玩家胜利
            else if(boardCase.winner===2&&outcome.innerText!="You lose!"){
                outcome.innerText="You lose!";
                alert("You lose!");
            }
        isAbled=false;
    }else if(outcome.innerText!=="Winner:") {

        outcome.innerText="Winner:";
    }
}



function initBoardCase(){
    boardCase={}
    boardCase.curPlayer=1;
    boardCase.winner=-1;
    boardCase.chesses=[];
    for(let i=0;i<10;i++){
        boardCase.chesses.push([])
        for(let j=0;j<10;j++){
            boardCase.chesses[i].push(0)
        }
    }
}


// 用来定时更新界面
function refreshUI(){
    // 根据棋盘数据修改界面
    for(let i=0;i<buttons.length;i++){
        const x=Math.floor(i/10);
        const y=i-10*x;
        // 获取棋子位置(x,y)后
        const kind=boardCase.chesses[x][y];
        if(kind===1&&buttons[i].innerHTML!=blackImg){
            buttons[i].innerHTML=blackImg;
        }else if(kind===2&&buttons[i].innerHTML!=whiteImg){
            buttons[i].innerHTML=whiteImg;
        }
        else if(kind===0&&(buttons[i].innerHTML!=backgroundImg)){
            buttons[i].innerHTML=backgroundImg;
        }
        else if(kind===-1&&buttons[i].innerHTML!=backgroundDeepImg){
            buttons[i].innerHTML=backgroundDeepImg;
        }
    }
}






// 读取内存,更新boardCase
function loadBoardCase(){
    s=sessionStorage.getItem("game");
    if(s===null){
        // 如果内存没有内容,则初始化buttons的内容
        // 每个位置用0/1/2分别表示棋盘，棋子1，棋子2
        initBoardCase();
        return;
    }
    // 否则转化为棋盘对象
    boardCase=JSON.parse(s);
    // TODO:检查是否传递过来的是正确的棋盘对象
    // if (boardCase['s']===null ){
    // }
}

// 往内存写入boardCase
function storeBoardCase(){
    s=JSON.stringify(boardCase);
    sessionStorage.setItem('game',s);
}