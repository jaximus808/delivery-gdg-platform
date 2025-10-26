
// import { auth, signOut } from "@/lib/auth" - 
// use nextAuth?

import Link from "next/link"

export default function Home() {
  return (
    <div className="min-h-screen">
      {/* Navigation */}
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <span className="text-2xl font-bold">Delivery Robot</span>
            </div>
            <div className="flex items-center gap-4">
              <a href="/" className="text-gray-700 font-medium">
                Log In
              </a>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left Content */}
          <div className="space-y-8">
            <h1 className="text-5xl sm:text-6xl font-bold text-gray-900 leading-tight">
              Autonomous robot delivery to your door
            </h1>
            <p className="text-xl text-gray-600">
              Order from your favorite campus restaurants and watch our robots deliver it right to you.
            </p>

            {/* Sign In */}
            <div className="bg-white rounded-2xl shadow-xl p-8 max-w-md">
              <h2 className="text-2xl font-semibold mb-6 text-gray-900">Get Started</h2>
              
              <Link href="/dashboard">
                <button className="w-full bg-green-400 hover:bg-green-300 text-white font-semibold py-3 px-6 rounded-lg transition-all duration-200">
                  Enter Dashboard (Demo)
                </button>
              </Link>
            </div>

            {/* Very real stats */}
            <div className="grid grid-cols-3 gap-6 pt-8">
              <div>
                <div className="text-3xl font-bold text-green-400">10000+</div>
                <div className="text-gray-600 text-sm">Active Robots</div>
              </div>
              <div>
                <div className="text-3xl font-bold text-green-400">25k+</div>
                <div className="text-gray-600 text-sm">Deliveries</div>
              </div>
              <div>
                <div className="text-3xl font-bold text-green-400">100%</div>
                <div className="text-gray-600 text-sm">Eco-Friendly</div>
              </div>
            </div>
          </div>

          {/* Resutraunt Examples */}
          <div className="hidden lg:block">
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-tr from-indigo-400 to-purple-400 rounded-3xl transform rotate-3"></div>
              <div className="relative bg-white rounded-3xl shadow-2xl p-8">
                <div className="space-y-4">
                  <div className="flex items-center gap-3 bg-indigo-50 p-4 rounded-lg">
                    <span className="text-2xl"></span>
                    <div className="flex-1">
                      <div className="font-semibold">Corner 17</div>
                      <div className="text-sm text-gray-600">Delivery Time • 15-25 min</div>
                    </div>
                    <span className="text-indigo-600 font-bold">
                      <img src='corner17_logo.png' alt="Corner 17" style={{ width: 80, height: 80, objectFit: "contain" }} />
                    </span>
                  </div>
                  <div className="flex items-center gap-3 bg-purple-50 p-4 rounded-lg">
                    <span className="text-2xl"></span>
                    <div className="flex-1">
                      <div className="font-semibold">Beast Craft Barbecue</div>
                      <div className="text-sm text-gray-600">Delivery Time • 20-30 min</div>
                    </div>
                    <span className="text-indigo-600 font-bold">
                      <img src='beast_craft_logo.png' alt="Beast Craft Barbecue" style={{ width: 80, height: 80, objectFit: "contain" }} />
                    </span>
                  </div>
                  <div className="flex items-center gap-3 bg-blue-50 p-4 rounded-lg">
                    <span className="text-2xl"></span>
                    <div className="flex-1">
                      <div className="font-semibold">QDOBA</div>
                      <div className="text-sm text-gray-600">Delivery Time • 10-20 min</div>
                    </div>
                    <span className="text-indigo-600 font-bold">
                      <img src='qdoba_logo.png' alt="QDOBA" style={{ width: 80, height: 80, objectFit: "contain" }} />
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>

      {/* Features */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="grid md:grid-cols-3 gap-8">
          <div className="text-center">
            <div className="mb-4 flex justify-center">
              <img src='robot_clipart.jpg' alt="Robot" style={{ width: 64, height: 64, objectFit: "contain" }} />
            </div>
            <h3 className="text-xl font-semibold mb-2">Autonomous Delivery</h3>
            <p className="text-gray-600">Watch your order arrive via our self-driving robots</p>
          </div>
          <div className="text-center">
            <div className="mb-4 flex justify-center">
              <img src='plant_clipart.png' alt="Plant" style={{ width: 64, height: 64, objectFit: "contain" }} />
            </div>
            <h3 className="text-xl font-semibold mb-2">Zero Emissions</h3>
            <p className="text-gray-600">100% electric, sustainable delivery</p>
          </div>
          <div className="text-center">
            <div className="mb-4 flex justify-center">
              <img src='pin_clipart.png' alt="Pin" style={{ width: 64, height: 64, objectFit: "contain" }} />
            </div>
            <h3 className="text-xl font-semibold mb-2">Real-Time Tracking</h3>
            <p className="text-gray-600">Follow your robot's journey to your door</p>
          </div>
        </div>
      </section>
    </div>
  )
}
