var testJson = [
  {
    "id": 1,
    "name": "金次郎"
  },
  {
    "id": 2,
    "name": "銀次郎"
  }
];


var App = Vue.component("App", {
    template: `

        <div class="container">
            <ul>
              <li v-for="item in objectItems1">{{item}}</li>

              <li v-for="(item, key, index) in objectItems2">
                {{ item }} - {{ key }} - {{ index }}
              </li>

              <li v-for="item in itemsLessThanTen" :key="item.name">
                {{ item.name }} - {{ item.price }}
              </li>

              <li v-for="item in Items3">{{item.id}}-{{item.name}}</li>
            </ul>

        </div>
    `,
    data(){
        return {
          objectItems1:{
            key1: 'item1',
            key2: 'item2',
            key3: 'item3'
          },
          objectItems2:{
            th:   'タイ語',
            yomi: '読み',
            imi:  '意味'
          },
          shoppingItems:[
            {name: 'apple', price: '7'},
            {name: 'orange', price: '12'},
            {name: 'banana', price: '11'}
          ],
          Items3:testJson,
          videoId: 'oy-HyrxsPJE',
          // hello: "Hello World!"};
        }
    },
    computed:{
      itemsLessThanTen: function() {
        return this.shoppingItems.filter(function(item) {
          return item.price > 10;
        })
      }
    }
});


new Vue({
    el: "#app"
});
