### Go API Showcase Project (Prototype)

โปรเจกต์นี้เป็นต้นแบบ (Prototype) ในการพัฒนา RESTful API ด้วยภาษา Go (Golang) เพื่อวัตถุประสงค์ในการจัดทำพอร์ตโฟลิโอสำหรับประกอบการพิจารณาสมัครงาน โดยเน้นการจัดโครงสร้างโค้ดที่อ่านง่าย แสดงถึงไสตล์การเขียน ขยายระบบได้ง่าย และรองรับการทำงานที่มีประสิทธิภาพ 

### 🚀 จุดประสงค์ของโปรเจกต์ (Purpose)

* **Portfolio Showcase:** แสดงทักษะการเขียนโค้ดภาษา Go และการออกแบบระบบ API ของผู้พัฒนา
* **Best Practices:** ประยุกต์ใช้โครงสร้างโฟลเดอร์แบบ Clean Architecture หรือ Standard Go Project Layout
* **Production-Ready Mindset:** จำลองสถาปัตยกรรมระบบที่มีความปลอดภัย และรองรับการทดสอบ (Testing)

### 🛠️ เทคโนโลยีที่ใช้ (Tech Stack)

* **Language:** Go 1.26+
* **Framework / Router:** Fiber
* **Database:** PostgreSQL พร้อม GORM (ORM)
* **Authentication:** JWT (JSON Web Token)
* **Testing:** Testify (`assert`, `mock`)
* **Other Tools:** Docker, Swagger (API Documentation)

### 🏗️ โครงสร้างสถาปัตยกรรม (Architecture)

โปรเจกต์นี้แยก Layer การทำงานอย่างชัดเจนเพื่อลด Decoupling (ความเกี่ยวเนื่องกันของโค้ด) 

* **Handler / Controller:** รับ Request และส่ง Response
* **Service / Usecase:** ประมวลผล Business Logic
* **Repository:** ติดต่อฐานข้อมูลและจัดการ Data Store

### ⏱️ วิธีการรันโปรเจกต์ (Getting Started)

1. Clone Repository นี้ลงเครื่องของคุณ :

    ```bash

    git clone https://github.com/markdeesoft/golang-api
    cd golang-api
    ```

2. ตั้งค่า Environment Variables (สร้างไฟล์ .env สามารถ copy จาก .env.example) :

    ```env

    PORT=8080
    ALLOW_ORIGIN="*"
    JWT_SECRET="secret"
    ...
    ```

3. การรัน Unit Test (Testing) :
    ```bash

    go test ./...
    ```

    * **`-cover`**: สรุปเปอร์เซ็นต์ Code Coverage ออกมาเป็นตัวเลขสั้นๆ ใน Terminal
        ```bash
        go test -cover ./...
        ```

    * **`-coverprofile` & `go tool cover`**: แสดงผลผ่าน Web Browser ซึ่งจะไฮไลต์ให้เห็นชัดเจนว่าบรรทัดไหนใน Handler หรือ Repository ที่ถูกรันแล้ว (สีเขียว) และบรรทัดไหนยังไม่ถูกทดสอบ (สีแดง)
        ```bash

        go test -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out
        ```


4. สั่งรันโปรเจกต์ผ่าน go run main.go :

    ```bash

    go run main.go
    ```

### 📝 API Endpoints (ตัวอย่าง)

* POST /auth/login - เข้าสู่ระบบเพื่อรับ JWT Token
* POST /admin/user - ลงทะเบียนผู้ใช้งานใหม่ (ต้องใส่ Bearer Token)
* GET /product - ดึงข้อมูลสินค้าและส่วนเชื่อมต่อข้อมูล (ต้องใส่ Bearer Token)
  
---

*หมายเหตุ: โปรเจกต์นี้สร้างขึ้นเพื่อจัดแสดงทักษะ (Showcase) หากต้องการนำไปใช้งานจริง จำเป็นต้องปรับปรุงการตั้งค่าความปลอดภัยให้รัดกุมยิ่งขึ้น*
