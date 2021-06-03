main :: IO ()
main = do
  i <- readLn
  putStrLn . show $ i * 10
