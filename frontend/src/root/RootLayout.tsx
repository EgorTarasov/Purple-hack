// import TopBar from "@/components/widgets/TopBar"
import { Outlet } from "react-router-dom"

const RootLayout = () => {
  return (
    <div className="w-full">
      {/* <TopBar /> */}
      {/* <LeftSideBar /> */}

      <section className="h-[100vh] mx-auto overflow-hidden">
        <Outlet />
      </section>

      {/* <BottomBar /> */}
    </div>
  )
}

export default RootLayout