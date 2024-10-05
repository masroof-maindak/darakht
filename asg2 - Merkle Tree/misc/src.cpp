#include <fstream>
#include <iostream>

// This program writes an arbitrary number of integers to a specified file in binary format

int main(int argc, char **argv) {
	if (argc != 3) {
		std::cerr << "Usage: " << argv[0]
				  << " <file_name> <intended_file_size (in MBs)>\n";
		return 1;
	}

	int fileSize		 = atoi(argv[2]);
	int numInts			 = (fileSize * 1024 * 1024) / sizeof(int);
	std::string fileName = argv[1];

	std::ofstream file(fileName, std::ios::out | std::ios::binary);

	if (file.is_open()) {
		file.write(reinterpret_cast<const char *>(&numInts), sizeof(numInts));

		for (int i = 0; i < numInts; i++) {
			int32_t randNo = rand() % 1000000;
			file.write(reinterpret_cast<const char *>(&randNo), sizeof(randNo));
		}

		file.close();
	} else {
		std::cerr << "Error opening file.\n";
		return 1;
	}

	return 0;
}
