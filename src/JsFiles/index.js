var todoCount = 0;
const nodeApi = "http://localhost/api/";

function addTask () {
    const taskName = document.getElementById("Todo-Input");
    const taskValue = taskName.value.trim();

    if (taskValue) {
        const div = document.createElement("div");
        const addTask = document.createElement("button");

        todoCount++;
        const todoNumber = String(todoCount);
        const todoId = "Todo-Number-" + todoNumber;

        div.className = "Todo-Lists-Add-Div";
        addTask.className = "Todo-Lists-Add-Item";
        addTask.textContent = taskValue;
        addTask.setAttribute("onclick", "todayToEndTaskMove(event)");
        addTask.id = todoId;

        div.appendChild(addTask);

        const todoListsAdd = document.getElementById("Todo-Lists-Add");
        todoListsAdd.appendChild(div);

        setTimeout(() => {
            div.classList.add("show"); // アニメーション用クラスを追加
          }, 10);

        taskName.value = "";
    }

    taskName.focus();
};

function todayToEndTaskMove (event) {
    const clickElement = event.target;
    const clickElementId = clickElement.id;

    const todayTasksDiv = document.getElementById("Todo-Lists-Add");
    const todayTasks = todayTasksDiv.querySelectorAll("div");

    const endTasksDiv = document.getElementById("Todo-Lists-End");
    const endTasks = endTasksDiv.querySelectorAll("div");

    todayTasks.forEach((todayTask) => {
        const todayTaskId = todayTask.querySelector("button").id;
        if (clickElementId === todayTaskId) {
            todayTask.className = "Todo-Lists-End-Div";
            endTasksDiv.appendChild(todayTask);

            setTimeout(() => {
                todayTask.classList.add("show"); // アニメーション用クラスを追加
              }, 5);
        }
    });

    endTasks.forEach ((endTask) => {
        const endTaskId = endTask.querySelector("button").id;
        if (clickElementId === endTaskId) {
            endTask.className = "Todo-Lists-Add-Div";
            todayTasksDiv.appendChild(endTask);

            setTimeout(() => {
                endTask.classList.add("show"); // アニメーション用クラスを追加
              }, 5);
        }
    });
};

document.addEventListener("DOMContentLoaded", () => {
    const inputElement = document.getElementById("Todo-Input");

    if (!inputElement) {
      console.error("Input element not found!");
      return;
    }

    inputElement.addEventListener("input", (event) => {
      const input = event.target;

      // 入力フィールドの位置とサイズを取得
      const rect = input.getBoundingClientRect();

      // キャレットの位置を取得
      const caretPosition = getCaretCoordinates(input);

      // 複数の星屑を生成
      for (let i = 0; i < 2; i++) {
        createSparkle(rect.left + caretPosition.x, rect.top + caretPosition.y);
      }
    });

    // 星屑を生成する関数
    function createSparkle(x, y) {
      const sparkle = document.createElement("span");
      sparkle.classList.add("sparkle");

      // ランダムな飛び散り方向
      const randomX = Math.random() * 50; // -100px ～ +100px
      const randomY = Math.random() * 50; // -100px ～ +100px
      const randomScale = Math.random() * 0.5 + 0.5; // 0.5 ～ 1.0のスケール
      const randomRotation = Math.random() * 360; // 回転角度

      // 星屑の初期位置
      sparkle.style.left = `${x}px`;
      sparkle.style.top = `${y}px`;

      // ランダムな変化を追加
      sparkle.style.setProperty("--x", `${randomX}px`);
      sparkle.style.setProperty("--y", `${randomY}px`);
      sparkle.style.setProperty("--scale", randomScale);
      sparkle.style.setProperty("--rotation", `${randomRotation}deg`);

      // 星屑をページに追加
      document.body.appendChild(sparkle);

      // アニメーション終了後に削除
      sparkle.addEventListener("animationend", () => {
        sparkle.remove();
      });
    }

    // キャレット位置を計算する関数
    function getCaretCoordinates(input) {
      const selectionStart = input.selectionStart;

      // 入力内容の前半部分を取得
      const text = input.value.substring(0, selectionStart);

      // 一時的なダミー要素を作成してキャレットの位置を計算
      const span = document.createElement("span");
      span.style.position = "absolute";
      span.style.visibility = "hidden";
      span.style.whiteSpace = "pre";
      span.style.font = window.getComputedStyle(input).font;
      span.textContent = text.replace(/\s/g, "\u00A0"); // スペースをノーブレークスペースに変換
      document.body.appendChild(span);

      const rect = span.getBoundingClientRect();
      const coordinates = { x: rect.width, y: rect.height / 2 }; // X座標を取得、高さの中心をYに
      span.remove();

      return coordinates;
    }
  });

function returnTaskArray() {
  const completedTasks = [];
  const incompleteTasks =[];

  const todayTasksDiv = document.getElementById("Todo-Lists-Add");
  const todayTasks = todayTasksDiv.querySelectorAll("div");
  const endTasksDiv = document.getElementById("Todo-Lists-End");
  const endTasks = endTasksDiv.querySelectorAll("div");

  todayTasks.forEach((todayTask) => {
    const todayTaskId = Number(todayTask.querySelector("button").id.match(/\d/));
    const todayTaskContent = todayTask.querySelector("button").textContent;
    const task = {
      "id": todayTaskId,
      "content": todayTaskContent
    };
    incompleteTasks.push(task);
  });

  endTasks.forEach((endTask) => {
    const endTaskId =  Number(endTask.querySelector("button").id.match(/\d/));
    const endTaskContent = endTask.querySelector("button").textContent;
    const task = {
      "id": endTaskId,
      "content": endTaskContent
    };
    completedTasks.push(task);
  });

  const taskArrays = {
    "completed":completedTasks,
    "incomplete": incompleteTasks
    };
  const taskJson = JSON.stringify(taskArrays);
  console.log(taskJson);
  const returnJson = sendJson(nodeApi, taskJson);
  return returnJson;
}

/**
 * 指定したURLへJSON文字列をPOSTリクエストする関数
 * @param {string} url - 送信先URL
 * @param {string} jsonData - JSON形式の文字列データ（すでにJSON文字列になっている前提）
 */
function sendJson(url, jsonData) {
  fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: jsonData // ここではstringifyしない
  })
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json();
  })
  .then(responseData => {
    console.log('Success:', responseData);
  })
  .catch(error => {
    console.error('Error:', error);
  });
}

// 使用例（jsonDataはすでにJSON文字列であることを想定）
// const dataString = '{"name":"John Doe","email":"john@example.com"}';
// sendJson('http://localhost:3000/api', dataString);